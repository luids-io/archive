// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package mongodb

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo"

	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/archive/pkg/archive/builder"
)

// Builder returns a builder function
func Builder() builder.BuildBackendFn {
	return func(b *builder.Builder, cfg builder.BackendDef) (archive.Backend, error) {
		if cfg.URL == "" {
			return nil, errors.New("'url' is required")
		}
		// create session with mgo
		session, err := mgo.Dial(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("dialing with mongodb '%s': %v", cfg.URL, err)
		}
		backend := &mdbBackend{session: session}
		b.OnShutdown(func() error {
			session.Close()
			return nil
		})
		return backend, nil
	}
}

const (
	// BackendClass defines backend name
	BackendClass = "mongodb"
)

func init() {
	builder.RegisterBackendBuilder(BackendClass, Builder())
}
