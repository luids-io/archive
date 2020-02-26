// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package backend

import (
	"errors"
	"fmt"

	"github.com/luids-io/archive/pkg/archive"
	"github.com/luisguillenc/yalogi"
)

// Builder constructs backends and services
type Builder struct {
	opts   options
	logger yalogi.Logger

	backends    map[string]bool
	backendList []*archive.Backend

	startup  []func() error
	shutdown []func() error
}

// Option is used for builder configuration
type Option func(*options)

type options struct {
	logger yalogi.Logger
}

var defaultOptions = options{logger: yalogi.LogNull}

// SetLogger sets a logger for the component
func SetLogger(l yalogi.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

// New instances a new builder
func New(opt ...Option) *Builder {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return &Builder{
		opts:        opts,
		logger:      opts.logger,
		backends:    make(map[string]bool),
		backendList: make([]*archive.Backend, 0),
		startup:     make([]func() error, 0),
		shutdown:    make([]func() error, 0),
	}
}

// Backend returns the Backend the id
func (b *Builder) Backend(id string) (*archive.Backend, bool) {
	for _, svc := range b.backendList {
		if svc.ID == id {
			return svc, true
		}
	}
	return nil, false
}

// Backends returns the backends created by builder
func (b *Builder) Backends() []*archive.Backend {
	backends := make([]*archive.Backend, len(b.backendList))
	copy(backends, b.backendList)
	return backends
}

// Build creates a Backend using the definition passed as param
func (b *Builder) Build(def Definition) (*archive.Backend, error) {
	b.logger.Debugf("building '%s' class '%s'", def.ID, def.Class)
	if def.ID == "" {
		return nil, errors.New("id field is required")
	}
	//check if exists
	_, ok := b.backends[def.ID]
	if ok {
		return nil, errors.New("'%s' exists")
	}
	//check if disabled
	if def.Disabled {
		return nil, fmt.Errorf("'%s' is disabled", def.ID)
	}
	//get builder
	customb, ok := registryBuilder[def.Class]
	if !ok {
		return nil, fmt.Errorf("can't find a builder for '%s' in '%s'", def.Class, def.ID)
	}
	n, err := customb(b, def) //builds
	if err != nil {
		return nil, fmt.Errorf("building '%s': %v", def.ID, err)
	}
	//register
	b.backends[def.ID] = true
	b.backendList = append(b.backendList, n)
	return n, nil
}

// OnStartup registers the functions that will be executed during startup.
func (b *Builder) OnStartup(f func() error) {
	b.startup = append(b.startup, f)
}

// OnShutdown registers the functions that will be executed during shutdown.
func (b *Builder) OnShutdown(f func() error) {
	b.shutdown = append(b.shutdown, f)
}

// Start executes all registered functions.
func (b *Builder) Start() error {
	b.logger.Infof("starting backend builder registered services")
	var ret error
	for _, f := range b.startup {
		err := f()
		if err != nil {
			return err
		}
	}
	return ret
}

// Shutdown executes all registered functions.
func (b *Builder) Shutdown() error {
	b.logger.Infof("shutting down backend builder registered services")
	var ret error
	for _, f := range b.shutdown {
		err := f()
		if err != nil {
			ret = err
		}
	}
	return ret
}

// Logger returns logger
func (b Builder) Logger() yalogi.Logger {
	return b.logger
}
