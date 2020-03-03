// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/luisguillenc/serverd"
	"github.com/luisguillenc/yalogi"
	"google.golang.org/grpc"

	dnsapi "github.com/luids-io/api/dnsutil/archive"
	eventapi "github.com/luids-io/api/event/archive"
	tlsapi "github.com/luids-io/api/tlsutil/archive"
	iconfig "github.com/luids-io/archive/internal/config"
	ifactory "github.com/luids-io/archive/internal/factory"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/archive/pkg/archive/backend"
	"github.com/luids-io/archive/pkg/archive/service"
	cconfig "github.com/luids-io/common/config"
	cfactory "github.com/luids-io/common/factory"

	// backends
	_ "github.com/luids-io/archive/pkg/archive/backend/mongodb"

	// services
	_ "github.com/luids-io/archive/pkg/archive/service/dnsmdb"
	_ "github.com/luids-io/archive/pkg/archive/service/eventmdb"
	_ "github.com/luids-io/archive/pkg/archive/service/tlsmdb"
)

func createLogger(debug bool) (yalogi.Logger, error) {
	cfgLog := cfg.Data("log").(*cconfig.LoggerCfg)
	return cfactory.Logger(cfgLog, debug)
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

func createBackends(srv *serverd.Manager, logger yalogi.Logger) (*backend.Builder, error) {
	cfgBackend := cfg.Data("backend").(*iconfig.BackendCfg)
	builder, err := ifactory.BackendBuilder(cfgBackend, logger)
	if err != nil {
		return nil, err
	}
	//create backends
	err = ifactory.Backends(cfgBackend, builder, logger)
	if err != nil {
		return nil, err
	}
	srv.Register(serverd.Service{
		Name:     "backend-builder.service",
		Start:    builder.Start,
		Shutdown: func() { builder.Shutdown() },
		Ping:     func() error { return builder.PingAll() },
	})
	return builder, nil
}

func createServices(srv *serverd.Manager, finder archive.BackendFinder, logger yalogi.Logger) (*service.Builder, error) {
	cfgService := cfg.Data("service").(*iconfig.ServiceCfg)
	builder, err := ifactory.ServiceBuilder(cfgService, finder, logger)
	if err != nil {
		return nil, err
	}
	//create services
	err = ifactory.Services(cfgService, builder, logger)
	if err != nil {
		return nil, err
	}
	srv.Register(serverd.Service{
		Name:     "service-builder.service",
		Start:    builder.Start,
		Shutdown: func() { builder.Shutdown() },
	})
	return builder, nil
}

func createEventAPIService(gsrv *grpc.Server, finder archive.ServiceFinder, logger yalogi.Logger) error {
	logger.Infof("creating and registering eventapi grpc service")
	cfgArchive := cfg.Data("archive").(*iconfig.ArchiveCfg)
	if cfgArchive.EventAPI != "" {
		gsvc, err := ifactory.EventAPIService(cfgArchive, finder, logger)
		if err != nil {
			return fmt.Errorf("creating eventapi service: %v", err)
		}
		eventapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}

func createDNSAPIService(gsrv *grpc.Server, finder archive.ServiceFinder, logger yalogi.Logger) error {
	logger.Infof("creating and registering dnsapi grpc service")
	cfgArchive := cfg.Data("archive").(*iconfig.ArchiveCfg)
	if cfgArchive.DNSAPI != "" {
		gsvc, err := ifactory.DNSAPIService(cfgArchive, finder, logger)
		if err != nil {
			return fmt.Errorf("creating dnsapi service: %v", err)
		}
		dnsapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}

func createTLSAPIService(gsrv *grpc.Server, finder archive.ServiceFinder, logger yalogi.Logger) error {
	logger.Infof("creating and registering tlsapi grpc service")
	cfgArchive := cfg.Data("archive").(*iconfig.ArchiveCfg)
	if cfgArchive.TLSAPI != "" {
		gsvc, err := ifactory.TLSAPIService(cfgArchive, finder, logger)
		if err != nil {
			return fmt.Errorf("creating tlsapi service: %v", err)
		}
		tlsapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}
