// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	"fmt"

	"github.com/globalsign/mgo"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/serverd"
	"github.com/luisguillenc/yalogi"

	iconfig "github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/eventarchive/mongodb"
	cconfig "github.com/luids-io/common/config"
	cfactory "github.com/luids-io/common/factory"
	"github.com/luids-io/core/event"
	"github.com/luids-io/core/event/services/archive"
)

func createLogger(debug bool) (yalogi.Logger, error) {
	cfgLog := cfg.Data("log").(*cconfig.LoggerCfg)
	return cfactory.Logger(cfgLog, debug)
}

// create tls archiver
func createArchiverSvc(srv *serverd.Manager, logger yalogi.Logger) (event.Archiver, error) {
	cfgArchiver := cfg.Data("").(*iconfig.ArchiverCfg)
	var backend event.Archiver

	switch cfgArchiver.Backend {
	case "mongodb":
		arc, err := createMDBArchiver(logger)
		if err != nil {
			return nil, fmt.Errorf("couldn't create backend: %v", err)
		}
		backend = arc
		srv.Register(serverd.Service{
			Name:     "mongodb.backend",
			Start:    arc.Start,
			Shutdown: arc.Shutdown,
			Ping:     arc.Ping,
		})
	default:
		return nil, fmt.Errorf("unknown backend '%s'", cfgArchiver.Backend)
	}
	return backend, nil
}

// create archiver server
func createArchiverSrv(a event.Archiver, srv *serverd.Manager, logger yalogi.Logger) error {
	//create server
	cfgServer := cfg.Data("grpc-archive").(*cconfig.ServerCfg)
	glis, gsrv, err := cfactory.Server(cfgServer)
	if err != nil {
		return err
	}
	// create service
	service := archive.NewService(a)
	archive.RegisterServer(gsrv, service)
	if cfgServer.Metrics {
		grpc_prometheus.Register(gsrv)
	}
	srv.Register(serverd.Service{
		Name:     "resolvarchive.server",
		Start:    func() error { go gsrv.Serve(glis); return nil },
		Shutdown: gsrv.GracefulStop,
		Stop:     gsrv.Stop,
	})
	return nil
}

func createHealthSrv(srv *serverd.Manager, logger yalogi.Logger) error {
	cfgHealth := cfg.Data("health").(*cconfig.HealthCfg)
	if !cfgHealth.Empty() {
		hlis, health, err := cfactory.Health(cfgHealth, srv, logger)
		if err != nil {
			logger.Fatalf("creating health server: %v", err)
		}
		srv.Register(serverd.Service{
			Name:     "health.server",
			Start:    func() error { go health.Serve(hlis); return nil },
			Shutdown: func() { health.Close() },
		})
	}
	return nil
}

// create mongodb archiver
func createMDBArchiver(logger yalogi.Logger) (*mongodb.Archiver, error) {
	cfgMongoDB := cfg.Data("mongodb").(*iconfig.MongoDBCfg)
	err := cfgMongoDB.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid mongodb: %v", err)
	}
	session, err := mgo.Dial(cfgMongoDB.URL)
	if err != nil {
		return nil, fmt.Errorf("dialing with mongodb '%s': %v", cfgMongoDB.URL, err)
	}
	archiver := mongodb.New(session, cfgMongoDB.Database,
		mongodb.SetLogger(logger),
		mongodb.SetPrefix(cfgMongoDB.Prefix))
	return archiver, nil
}
