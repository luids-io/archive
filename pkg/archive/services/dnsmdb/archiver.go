// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package dnsmdb implements dnsutil.Archive using mongodb backend.
//
// This package is a work in progress and makes no API stability promises.
package dnsmdb

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/luids-io/api/dnsutil"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/archive/pkg/mongoutil"
	"github.com/luids-io/core/yalogi"
)

// ServiceClass registered.
const ServiceClass = "dnsmdb"

// Collection names.
const (
	ResolvColName = "resolvs"
)

// Default values.
const (
	DefaultResolvBulkSize = 1024
	DefaultSyncSeconds    = 5
)

// Archiver implements dns archive backend using a mongo database.
type Archiver struct {
	id     string
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
	bulkResolvs *mongoutil.Bulk
}

// New creates a new storage.
func New(id string, session *mgo.Session, db string, opt ...Option) *Archiver {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	s := &Archiver{
		id:       id,
		opts:     opts,
		logger:   opts.logger,
		database: db,
		session:  session,
	}
	return s
}

// Option encapsules options.
type Option func(*options)

type options struct {
	logger         yalogi.Logger
	closeSession   bool
	resolvBulkSize int
	syncSecs       int
	prefix         string
}

var defaultOptions = options{
	logger:         yalogi.LogNull,
	resolvBulkSize: DefaultResolvBulkSize,
	syncSecs:       DefaultSyncSeconds,
	closeSession:   false,
}

// SetLogger option allows set a custom logger.
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// CloseSession option allows close mongo session on shutdown.
func CloseSession(b bool) Option {
	return func(o *options) {
		o.closeSession = b
	}
}

// SetPrefix option allows set a prefix to collection.
func SetPrefix(s string) Option {
	return func(o *options) {
		o.prefix = s
	}
}

// Start the archiver.
func (a *Archiver) Start() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.started {
		return fmt.Errorf("archiver started")
	}
	a.logger.Infof("%s: starting mongodb dns archiver", a.id)
	//create indexes
	err := a.createIdx()
	if err != nil {
		return err
	}
	//init bulks & caches
	a.bulkResolvs = mongoutil.NewBulk(
		a.getCollection(ResolvColName),
		a.opts.resolvBulkSize,
	)
	//init control
	a.close = make(chan struct{})
	go a.doSync()
	a.started = true
	return nil
}

// SaveResolv implements dnsutil.Archiver interface.
func (a *Archiver) SaveResolv(ctx context.Context, r *dnsutil.ResolvData) (string, error) {
	if !a.started {
		return "", dnsutil.ErrUnavailable
	}
	r.ID = bson.NewObjectId().String()
	err := a.bulkResolvs.Insert(toMongoData(r))
	if err != nil {
		a.logger.Warnf("%s: saving resolv '%s': %v", a.id, r.Name, err)
		return "", dnsutil.ErrInternal
	}
	return r.ID, nil
}

// Shutdown closes the conection.
func (a *Archiver) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.started {
		a.logger.Infof("%s: shutting down dns archiver", a.id)
		a.started = false
		close(a.close)
		a.session.Fsync(false)
		if a.opts.closeSession {
			a.session.Close()
		}
	}
	return
}

// Ping tests the connection with the storage.
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
				a.logger.Warnf("%s: %v", a.id, err)
			}
		case <-a.close:
			errs := a.syncBulks()
			for _, err := range errs {
				a.logger.Warnf("%s: %v", a.id, err)
			}
			break
		}
	}
}

func (a *Archiver) syncBulks() []error {
	errs := make([]error, 0, 1)
	var err error
	err = a.bulkResolvs.Flush()
	if err != nil {
		errs = append(errs, fmt.Errorf("sync resolvs: %v", err))
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

// ID implements archive.Service interface.
func (a *Archiver) ID() string {
	return a.id
}

// Class implements archive.Service interface.
func (a *Archiver) Class() string {
	return ServiceClass
}

// Implements implements archive.Service interface.
func (a *Archiver) Implements() []archive.API {
	return []archive.API{archive.DNSAPI}
}
