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
	"github.com/google/uuid"
	"golang.org/x/net/publicsuffix"

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
	DefaultDBName         = "luidsdb"
	DefaultResolvBulkSize = 1024
	DefaultSyncSeconds    = 5
	DefaultMaxSize        = 100
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
func (a *Archiver) SaveResolv(ctx context.Context, rd dnsutil.ResolvData) (uuid.UUID, error) {
	if !a.started {
		return uuid.Nil, dnsutil.ErrUnavailable
	}
	// create new uuid if not set
	sid := rd.ID.String()
	if sid == "" {
		newid, err := uuid.NewRandom()
		if err != nil {
			a.logger.Warnf("%s: saveresolv(): generating new id: %v", a.id, err)
			return uuid.Nil, dnsutil.ErrInternal
		}
		rd.ID = newid
		sid = rd.ID.String()
	}
	// compute fields
	rd.TLD, _ = publicsuffix.PublicSuffix(rd.Name)
	rd.TLDPlusOne, _ = publicsuffix.EffectiveTLDPlusOne(rd.Name)
	// convert to mongo data
	m := &mdbResolvData{}
	err := toMData(&rd, m)
	if err != nil {
		a.logger.Warnf("%s: saveresolv(%s): converting to mongo: %v", a.id, sid, err)
		return uuid.Nil, dnsutil.ErrBadRequest
	}
	// store data
	m.StorageID = bson.NewObjectId()
	err = a.bulkResolvs.Insert(m)
	if err != nil {
		a.logger.Warnf("%s: saveresolv(%s): inserting in bulk: %v", a.id, sid, err)
		return uuid.Nil, dnsutil.ErrInternal
	}
	return rd.ID, nil
}

// GetResolv implements dnsutil.Finder interface.
func (a *Archiver) GetResolv(ctx context.Context, id uuid.UUID) (dnsutil.ResolvData, bool, error) {
	if !a.started {
		return dnsutil.ResolvData{}, false, dnsutil.ErrUnavailable
	}
	//if invalid id, then returns not found
	sid := id.String()
	if sid == "" {
		return dnsutil.ResolvData{}, false, nil
	}
	//do find
	var m mdbResolvData
	c := a.getCollection(ResolvColName)
	err := c.Find(bson.M{"id": sid}).One(&m)
	if err == mgo.ErrNotFound {
		return dnsutil.ResolvData{}, false, nil
	}
	if err != nil {
		a.logger.Warnf("%s: getresolv(%s): %v", a.id, sid, err)
		return dnsutil.ResolvData{}, false, dnsutil.ErrInternal
	}
	//encode response
	var r dnsutil.ResolvData
	err = fromMData(&m, &r)
	if err != nil {
		a.logger.Warnf("%s: getresolv(%s): converting from mongo: %v", a.id, sid, err)
		return dnsutil.ResolvData{}, false, dnsutil.ErrInternal
	}
	return r, true, nil
}

// ListResolvs implements dnsutil.Finder interface.
func (a *Archiver) ListResolvs(ctx context.Context, filters []dnsutil.ResolvsFilter,
	rev bool, max int, next string) ([]dnsutil.ResolvData, string, error) {
	if !a.started {
		return nil, "", dnsutil.ErrUnavailable
	}
	c := a.getCollection(ResolvColName)
	//create filter
	filter := createFilter(filters)
	if next != "" && bson.IsObjectIdHex(next) {
		if rev {
			filter["_id"] = bson.M{"$lt": bson.ObjectIdHex(next)}
		} else {
			filter["_id"] = bson.M{"$gt": bson.ObjectIdHex(next)}
		}
	}
	//do find
	q := c.Find(filter)
	if rev {
		q = q.Sort("-_id")
	}
	if max == 0 && DefaultMaxSize > 0 {
		max = DefaultMaxSize
	}
	if max > 0 {
		q = q.Limit(max)
	}
	//do query
	var mdbAll []mdbResolvData
	err := q.All(&mdbAll)
	if err != nil {
		a.logger.Warnf("%s: listresolvs(): %v", a.id, err)
		return nil, "", dnsutil.ErrInternal
	}
	//convert data
	last := ""
	result := make([]dnsutil.ResolvData, 0, len(mdbAll))
	for _, m := range mdbAll {
		var r dnsutil.ResolvData
		err = fromMData(&m, &r)
		if err != nil {
			a.logger.Warnf("%s: listresolvs(): converting from mongo '%s': %v", a.id, m.StorageID.Hex(), err)
			return nil, "", dnsutil.ErrInternal
		}
		result = append(result, r)
		last = m.StorageID.Hex()
	}
	//return
	if max > 0 && len(result) == max {
		return result, last, nil
	}
	return result, "", nil
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

func createFilter(filters []dnsutil.ResolvsFilter) bson.M {
	switch len(filters) {
	case 0:
		return bson.M{}
	case 1:
		return bsonFilter(filters[0])
	}
	mfilters := make([]bson.M, 0, len(filters))
	for _, f := range filters {
		mfilters = append(mfilters, bsonFilter(f))
	}
	return bson.M{"$or": mfilters}
}

func bsonFilter(f dnsutil.ResolvsFilter) bson.M {
	m := make(bson.M)
	if !f.Since.IsZero() || !f.To.IsZero() {
		tfilter := bson.M{}
		if !f.Since.IsZero() {
			tfilter["$gt"] = f.Since
		}
		if !f.To.IsZero() {
			tfilter["$lt"] = f.To
		}
		m["timestamp"] = tfilter
	}
	if f.Client != nil {
		m["clientIP"] = f.Client.String()
	}
	if f.Server != nil {
		m["serverIP"] = f.Server.String()
	}
	if f.Name != "" {
		m["name"] = f.Name
	}
	if f.ResolvedIP != nil {
		m["resolvedIPs"] = f.ResolvedIP.String()
	}
	if f.ResolvedCNAME != "" {
		m["resolvedCNAMEs"] = f.ResolvedCNAME
	}
	if f.QID > 0 {
		m["qid"] = f.QID
	}
	if f.ReturnCode > 0 {
		m["returnCode"] = f.ReturnCode
	}
	if f.TLD != "" {
		m["tld"] = f.TLD
	}
	if f.TLDPlusOne != "" {
		m["tldPlusOne"] = f.TLDPlusOne
	}
	return m
}
