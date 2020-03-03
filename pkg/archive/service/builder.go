// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package service

import (
	"errors"
	"fmt"

	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/archive/pkg/archive"
)

// Builder constructs archive services from definitions
type Builder struct {
	archive.ServiceFinder

	opts     options
	logger   yalogi.Logger
	services map[string]archive.Service
	bfinder  archive.BackendFinder
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

// NewBuilder instances a new builder
func NewBuilder(finder archive.BackendFinder, opt ...Option) *Builder {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return &Builder{
		opts:     opts,
		logger:   opts.logger,
		services: make(map[string]archive.Service),
		bfinder:  finder,
		startup:  make([]func() error, 0),
		shutdown: make([]func() error, 0),
	}
}

// FindServiceByID returns the Service with the id
func (b *Builder) FindServiceByID(id string) (archive.Service, bool) {
	svc, ok := b.services[id]
	return svc, ok
}

// Build creates a Service using the definition passed as param
func (b *Builder) Build(def Definition) (archive.Service, error) {
	b.logger.Debugf("building '%s' class '%s'", def.ID, def.Class)
	if def.ID == "" {
		return nil, errors.New("id field is required")
	}
	//check if exists
	_, ok := b.services[def.ID]
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
	b.services[def.ID] = n
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
	b.logger.Infof("starting service builder services")
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
	b.logger.Infof("shutting down service builder services")
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

// BackendFinder returns backend Finder
func (b Builder) BackendFinder() archive.BackendFinder {
	return b.bfinder
}
