// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package archive

import (
	"errors"
	"fmt"
	"strings"

	"github.com/luids-io/core/yalogi"
)

// Builder constructs backends and services from definitions
type Builder struct {
	opts     options
	logger   yalogi.Logger
	backends map[string]Backend
	services map[string]Service
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
func NewBuilder(opt ...Option) *Builder {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return &Builder{
		opts:     opts,
		logger:   opts.logger,
		backends: make(map[string]Backend),
		services: make(map[string]Service),
		startup:  make([]func() error, 0),
		shutdown: make([]func() error, 0),
	}
}

// Backend returns the Backend with the id
func (b *Builder) Backend(id string) (Backend, bool) {
	ba, ok := b.backends[id]
	return ba, ok
}

// Service returns the Service with the id
func (b *Builder) Service(id string) (Service, bool) {
	svc, ok := b.services[id]
	return svc, ok
}

// BuildBackend creates a Backend using the definition passed as param
func (b *Builder) BuildBackend(def BackendDef) (Backend, error) {
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
	customb, ok := regBackendBuilder[def.Class]
	if !ok {
		return nil, fmt.Errorf("can't find a builder for '%s' in '%s'", def.Class, def.ID)
	}
	n, err := customb(b, def) //builds
	if err != nil {
		return nil, fmt.Errorf("building '%s': %v", def.ID, err)
	}
	//register
	b.backends[def.ID] = n
	return n, nil
}

// BuildService creates a Service using the definition passed as param
func (b *Builder) BuildService(def ServiceDef) (Service, error) {
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
	customb, ok := regServiceBuilder[def.Class]
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
	b.logger.Infof("starting archive services")
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
	b.logger.Infof("shutting down archive services")
	errs := make([]string, 0, len(b.shutdown))
	for i := len(b.shutdown) - 1; i >= 0; i-- {
		fn := b.shutdown[i]
		err := fn()
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}
	return nil
}

// PingAll backends.
func (b *Builder) PingAll() error {
	b.logger.Debugf("PingAll()")
	errs := make([]string, 0, len(b.backends))
	for k, v := range b.backends {
		err := v.Ping()
		if err != nil {
			errs = append(errs, fmt.Sprintf("backend '%s': %v", k, err))
		}
	}
	if len(errs) > 0 {
		retErr := errors.New(strings.Join(errs, ";"))
		b.logger.Warnf("%s", retErr)
		return retErr
	}
	return nil
}

// Logger returns logger
func (b Builder) Logger() yalogi.Logger {
	return b.logger
}
