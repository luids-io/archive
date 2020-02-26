// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

package tlsmdb

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/luisguillenc/mongoutil"
	"github.com/luisguillenc/yalogi"
	cache "github.com/patrickmn/go-cache"

	"github.com/luids-io/core/tlsutil"
)

// Collection names
const (
	ConnectionColName  = "connections"
	CertificateColName = "certificates"
	RecordsColName     = "records"
)

// Default values
const (
	DefaultConnsBulkSize        = 256
	DefaultRecordsBulkSize      = 1024
	DefaultSyncSeconds          = 5
	DefaultCacheCertsExpiration = 30 * time.Minute
	DefaultCacheCertsCleanUp    = 5 * time.Minute
)

// Archiver implements resolv archive backend using a mongo database
type Archiver struct {
	opts   options
	logger yalogi.Logger
	//database
	session  *mgo.Session
	database string
	//control
	mu      sync.Mutex
	started bool
	close   chan struct{}
	//bulks & caches
	bulkConns   *mongoutil.Bulk
	bulkRecords *mongoutil.Bulk
	cacheCerts  *cache.Cache
}

// New creates a new mongodb storage
func New(session *mgo.Session, db string, opt ...Option) *Archiver {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	s := &Archiver{
		opts:     opts,
		logger:   opts.logger,
		database: db,
		session:  session,
	}
	return s
}

// Option encapsules options
type Option func(*options)

type options struct {
	logger               yalogi.Logger
	connsBulkSize        int
	recordsBulkSize      int
	syncSecs             int
	cacheCertsExpiration time.Duration
	cacheCertsCleanUp    time.Duration
	closeSession         bool
	prefix               string
}

var defaultOptions = options{
	logger:               yalogi.LogNull,
	connsBulkSize:        DefaultConnsBulkSize,
	recordsBulkSize:      DefaultRecordsBulkSize,
	syncSecs:             DefaultSyncSeconds,
	cacheCertsExpiration: DefaultCacheCertsExpiration,
	cacheCertsCleanUp:    DefaultCacheCertsCleanUp,
}

// SetLogger option allows set a custom logger
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// CloseSession option allows close mongo session on shutdown
func CloseSession(b bool) Option {
	return func(o *options) {
		o.closeSession = b
	}
}

// SetPrefix option allows set a prefix to collection
func SetPrefix(s string) Option {
	return func(o *options) {
		o.prefix = s
	}
}

// Start the archiver
func (a *Archiver) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.started {
		return fmt.Errorf("archiver started")
	}
	a.logger.Infof("starting mongodb tls archiver")
	//create indexes
	err := a.createIdx()
	if err != nil {
		return err
	}
	//init bulks & caches
	a.bulkConns = mongoutil.NewBulk(
		a.getCollection(ConnectionColName),
		a.opts.connsBulkSize,
	)
	a.bulkRecords = mongoutil.NewBulk(
		a.getCollection(RecordsColName),
		a.opts.recordsBulkSize,
	)
	a.cacheCerts = cache.New(
		DefaultCacheCertsExpiration,
		DefaultCacheCertsCleanUp,
	)
	//init control
	a.close = make(chan struct{})
	go a.doSync()
	a.started = true
	return nil
}

// SaveConnection implements tlsutil.Archiver interface
func (a *Archiver) SaveConnection(ctx context.Context, cn *tlsutil.ConnectionData) (string, error) {
	if !a.started {
		return "", fmt.Errorf("archiver not started")
	}
	return cn.ID, a.bulkConns.Insert(cn)
}

// SaveCertificate implements tlsutil.Archiver interface
func (a *Archiver) SaveCertificate(ctx context.Context, cert *tlsutil.CertificateData) (string, error) {
	if !a.started {
		return "", fmt.Errorf("archiver not started")
	}
	// check in cache
	ccert, ok := a.cacheCerts.Get(cert.Digest)
	if ok {
		cert, _ = ccert.(*tlsutil.CertificateData)
		return cert.ID, nil
	}
	// check in database
	var dbcert tlsutil.CertificateData
	err := a.getCollection(CertificateColName).
		Find(bson.M{"digest": cert.Digest}).One(&dbcert)
	if err != nil && err != mgo.ErrNotFound {
		a.logger.Errorf("finding cert digest: %v", err)
		return "", err
	} else if err == nil {
		//exists, but not in cache-> add to cache
		a.cacheCerts.Add(cert.Digest, &dbcert, cache.DefaultExpiration)
		return dbcert.ID, nil
	}
	// don't exist, add to cache
	a.cacheCerts.Add(cert.Digest, cert, cache.DefaultExpiration)
	return cert.ID, a.getCollection(CertificateColName).Insert(cert)
}

// StoreRecord implements tlsutil.Archiver interface
func (a *Archiver) StoreRecord(r *tlsutil.RecordData) error {
	if !a.started {
		return fmt.Errorf("archiver not started")
	}
	return a.bulkRecords.Insert(r)
}

// Shutdown closes the conection
func (a *Archiver) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.started {
		a.logger.Infof("shutting down tls archiver")
		a.started = false
		close(a.close)
		a.session.Fsync(false)
		if a.opts.closeSession {
			a.session.Close()
		}
	}
	return
}

// Ping tests the connection with the storage
func (a *Archiver) Ping() error {
	a.logger.Debugf("ping")
	if !a.started {
		return errors.New("archiver not started")
	}
	return a.session.Ping()
}

func (a *Archiver) doSync() {
	tick := time.NewTicker(time.Duration(a.opts.syncSecs) * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			errs := a.syncBulks()
			for _, err := range errs {
				a.logger.Warnf("%v", err)
			}
		case <-a.close:
			errs := a.syncBulks()
			for _, err := range errs {
				a.logger.Warnf("%v", err)
			}
			break
		}
	}
}

func (a *Archiver) syncBulks() []error {
	errs := make([]error, 0, 2)
	var err error
	err = a.bulkConns.Flush()
	if err != nil {
		errs = append(errs, fmt.Errorf("sync connections: %v", err))
	}
	err = a.bulkRecords.Flush()
	if err != nil {
		errs = append(errs, fmt.Errorf("sync records: %v", err))
	}
	return errs
}

func (a *Archiver) getDatabase() *mgo.Database {
	return a.session.DB(a.database)
}

func (a *Archiver) getCollection(name string) *mgo.Collection {
	if a.opts.prefix != "" {
		name = a.opts.prefix + "_" + name
	}
	return a.session.DB(a.database).C(name)
}
