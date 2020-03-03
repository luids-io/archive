// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package eventmdb

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo"

	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/archive/pkg/archive/backend/mongodb"
	"github.com/luids-io/archive/pkg/archive/service"
	"github.com/luids-io/core/option"
)

// Builder returns a builder function
func Builder() service.BuildFn {
	return func(b *service.Builder, cfg service.Definition) (*archive.Service, error) {
		if cfg.Backend == "" {
			return nil, errors.New("'backend' is required")
		}
		//get mongodb backend
		back, ok := b.BackendFinder().FindBackendByID(cfg.Backend)
		if !ok {
			return nil, errors.New("'backend' not found")
		}
		if back.GetClass() != mongodb.BackendClass {
			return nil, fmt.Errorf("'backend' class '%s' not suported in service", back.GetClass())
		}
		// get session from backend container
		session, ok := back.GetSession().(*mgo.Session)
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
		//create service container
		svc := &archive.Service{
			ID:     cfg.ID,
			Class:  ServiceClass,
			API:    archive.EventAPI,
			Object: archiver,
		}
		return svc, nil
	}
}

const (
	// ServiceClass defines service name
	ServiceClass = "eventmdb"
)

func init() {
	service.RegisterBuilder(ServiceClass, Builder())
}
