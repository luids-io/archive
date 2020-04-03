// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package eventmdb

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo"

	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/archive/pkg/archive/backends/mongodb"
	"github.com/luids-io/archive/pkg/archive/builder"
	"github.com/luids-io/core/utils/option"
)

// Builder returns a builder function
func Builder() builder.BuildServiceFn {
	return func(b *builder.Builder, cfg builder.ServiceDef) (archive.Service, error) {
		if cfg.Backend == "" {
			return nil, errors.New("'backend' is required")
		}
		//get mongodb backend
		back, ok := b.Backend(cfg.Backend)
		if !ok {
			return nil, errors.New("'backend' not found")
		}
		if back.Class() != mongodb.BackendClass {
			return nil, fmt.Errorf("'backend' class '%s' not suported in service", back.Class())
		}
		// get session from backend container
		session, ok := back.Session().(*mgo.Session)
		if !ok {
			return nil, errors.New("'backend' not found")
		}
		// parse options
		bopt := make([]Option, 0)
		bopt = append(bopt, SetLogger(b.Logger()))
		//by default, it uses id as database name
		dbname := cfg.ID
		if cfg.Opts != nil {
			var err error
			dbnameOpt, ok, err := option.String(cfg.Opts, "dbname")
			if err != nil {
				return nil, err
			}
			if ok {
				dbname = dbnameOpt
			}
			prefixOpt, ok, err := option.String(cfg.Opts, "prefix")
			if err != nil {
				return nil, err
			}
			if ok {
				bopt = append(bopt, SetPrefix(prefixOpt))
			}
		}
		//create archive service
		archiver := New(session, dbname, bopt...)
		b.OnStartup(func() error {
			return archiver.Start()
		})
		b.OnShutdown(func() error {
			archiver.Shutdown()
			return nil
		})
		return archiver, nil
	}
}

const (
	// ServiceClass defines service name
	ServiceClass = "eventmdb"
)

func init() {
	builder.RegisterServiceBuilder(ServiceClass, Builder())
}