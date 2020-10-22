// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. See LICENSE.

// Package eventmdb implements event.Archive using mongodb backend.
//
// This package is a work in progress and makes no API stability promises.
package eventmdb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/globalsign/mgo"

	"github.com/luids-io/api/event"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/core/yalogi"
)

// ServiceClass registered.
const ServiceClass = "eventmdb"

// Collection names.
const (
	EventColName = "events"
)

// Default values.
const (
	DefaultDBName = "luidsdb"
)

// Archiver implements event archive backend using a mongo database.
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
	logger       yalogi.Logger
	closeSession bool
	prefix       string
}

var defaultOptions = options{logger: yalogi.LogNull}

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
	a.logger.Infof("%s: starting mongodb event archiver", a.id)
	//create indexes
	err := a.createIdx()
	if err != nil {
		return err
	}
	//init control
	a.started = true
	return nil
}

// SaveEvent implements event.Archiver interface.
func (a *Archiver) SaveEvent(ctx context.Context, e event.Event) (string, error) {
	if !a.started {
		return "", event.ErrUnavailable
	}
	err := a.getCollection(EventColName).Insert(e)
	if err != nil {
		a.logger.Warnf("%s: saving event '%s': %v", a.id, e.ID, err)
		return "", event.ErrInternal
	}
	return e.ID, nil
}

// Shutdown closes the conection.
func (a *Archiver) Shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.started {
		a.logger.Infof("%s: shutting down event archiver", a.id)
		a.started = false
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

// Class implements archive.Service interface
func (a *Archiver) Class() string {
	return ServiceClass
}

// Implements implements archive.Service interface.
func (a *Archiver) Implements() []archive.API {
	return []archive.API{archive.EventAPI}
}
