// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/serverd"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"

	iconfig "github.com/luids-io/archive/internal/config"
	ifactory "github.com/luids-io/archive/internal/factory"
	cconfig "github.com/luids-io/common/config"
	cfactory "github.com/luids-io/common/factory"
	"github.com/luids-io/core/dnsutil"
	dnsarchive "github.com/luids-io/core/dnsutil/services/archive"
	"github.com/luids-io/core/event"
	eventarchive "github.com/luids-io/core/event/services/archive"
	tlsarchive "github.com/luids-io/core/tlsutil/services/archive"
	"github.com/luids-io/core/tlsutil"
)

func createLogger(debug bool) (yalogi.Logger, error) {
	cfgLog := cfg.Data("log").(*cconfig.LoggerCfg)
	return cfactory.Logger(cfgLog, debug)
}

// create archiver services
func createArchiverSvcs(gsrv *grpc.Server, srv *serverd.Manager, logger yalogi.Logger) error {
	cfgArchiver := cfg.Data("").(*iconfig.ArchiverCfg)
	if hasString(cfgArchiver.Services, "dns") {
		a, err := createDNSArchiverSvc(srv, logger)
		if err != nil {
			return err
		}
		service := dnsarchive.NewService(a)
		dnsarchive.RegisterServer(gsrv, service)
	}
	if hasString(cfgArchiver.Services, "event") {
		a, err := createEventArchiverSvc(srv, logger)
		if err != nil {
			return err
		}
		service := eventarchive.NewService(a)
		eventarchive.RegisterServer(gsrv, service)
	}
	if hasString(cfgArchiver.Services, "tls") {
		a, err := createTLSArchiverSvc(srv, logger)
		if err != nil {
			return err
		}
		service := tlsarchive.NewService(a)
		tlsarchive.RegisterServer(gsrv, service)
	}
	return nil
}

// create dns archiver
func createDNSArchiverSvc(srv *serverd.Manager, logger yalogi.Logger) (dnsutil.Archiver, error) {
	cfgArchiver := cfg.Data("").(*iconfig.ArchiverCfg)
	var backend dnsutil.Archiver

	switch cfgArchiver.Backend {
	case "mongodb":
		cfgMongoDB := cfg.Data("mongodb").(*iconfig.MongoDBCfg)
		arc, err := ifactory.DNSArchiveMDB(cfgMongoDB, logger)
		if err != nil {
			return nil, fmt.Errorf("couldn't create dns-archive backend: %v", err)
		}
		backend = arc
		srv.Register(serverd.Service{
			Name:     "dns-archive.backend",
			Start:    arc.Start,
			Shutdown: arc.Shutdown,
			Ping:     arc.Ping,
		})
	default:
		return nil, fmt.Errorf("unknown backend '%s'", cfgArchiver.Backend)
	}
	return backend, nil
}

// create event archiver
func createEventArchiverSvc(srv *serverd.Manager, logger yalogi.Logger) (event.Archiver, error) {
	cfgArchiver := cfg.Data("").(*iconfig.ArchiverCfg)
	var backend event.Archiver

	switch cfgArchiver.Backend {
	case "mongodb":
		cfgMongoDB := cfg.Data("mongodb").(*iconfig.MongoDBCfg)
		arc, err := ifactory.EventArchiveMDB(cfgMongoDB, logger)
		if err != nil {
			return nil, fmt.Errorf("couldn't create dns-archive backend: %v", err)
		}
		backend = arc
		srv.Register(serverd.Service{
			Name:     "event-archive.backend",
			Start:    arc.Start,
			Shutdown: arc.Shutdown,
			Ping:     arc.Ping,
		})
	default:
		return nil, fmt.Errorf("unknown backend '%s'", cfgArchiver.Backend)
	}
	return backend, nil
}

// create tls archiver
func createTLSArchiverSvc(srv *serverd.Manager, logger yalogi.Logger) (tlsutil.Archiver, error) {
	cfgArchiver := cfg.Data("").(*iconfig.ArchiverCfg)
	var backend tlsutil.Archiver

	switch cfgArchiver.Backend {
	case "mongodb":
		cfgMongoDB := cfg.Data("mongodb").(*iconfig.MongoDBCfg)
		arc, err := ifactory.TLSArchiveMDB(cfgMongoDB, logger)
		if err != nil {
			return nil, fmt.Errorf("couldn't create dns-archive backend: %v", err)
		}
		backend = arc
		srv.Register(serverd.Service{
			Name:     "tls-archive.backend",
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
func createArchiverSrv(srv *serverd.Manager, logger yalogi.Logger) (*grpc.Server, error) {
	//create server
	cfgServer := cfg.Data("grpc-archive").(*cconfig.ServerCfg)
	glis, gsrv, err := cfactory.Server(cfgServer)
	if err != nil {
		return nil, err
	}
	if cfgServer.Metrics {
		grpc_prometheus.Register(gsrv)
	}
	srv.Register(serverd.Service{
		Name:     "grpc-archive.server",
		Start:    func() error { go gsrv.Serve(glis); return nil },
		Shutdown: gsrv.GracefulStop,
		Stop:     gsrv.Stop,
	})
	return gsrv, nil
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

func hasString(ss []string, s string) bool {
	for _, b := range ss {
		if b == s {
			return true
		}
	}
	return false
}
