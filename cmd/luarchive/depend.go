// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
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
	"github.com/luids-io/core/utils/serverd"
	"github.com/luids-io/core/utils/yalogi"
)

func createLogger(debug bool) (yalogi.Logger, error) {
	cfgLog := cfg.Data("log").(*cconfig.LoggerCfg)
	return cfactory.Logger(cfgLog, debug)
}

func createHealthSrv(msrv *serverd.Manager, logger yalogi.Logger) error {
	cfgHealth := cfg.Data("health").(*cconfig.HealthCfg)
	if !cfgHealth.Empty() {
		hlis, health, err := cfactory.Health(cfgHealth, msrv, logger)
		if err != nil {
			logger.Fatalf("creating health server: %v", err)
		}
		msrv.Register(serverd.Service{
			Name:     "health.server",
			Start:    func() error { go health.Serve(hlis); return nil },
			Shutdown: func() { health.Close() },
		})
	}
	return nil
}

func createArchiverSrv(msrv *serverd.Manager) (*grpc.Server, error) {
	//create server
	cfgServer := cfg.Data("server-archive").(*cconfig.ServerCfg)
	glis, gsrv, err := cfactory.Server(cfgServer)
	if err != nil {
		return nil, err
	}
	if cfgServer.Metrics {
		grpc_prometheus.Register(gsrv)
	}
	msrv.Register(serverd.Service{
		Name:     fmt.Sprintf("[%s].server", cfgServer.ListenURI),
		Start:    func() error { go gsrv.Serve(glis); return nil },
		Shutdown: gsrv.GracefulStop,
		Stop:     gsrv.Stop,
	})
	return gsrv, nil
}

func createBackends(msrv *serverd.Manager, logger yalogi.Logger) (*backend.Builder, error) {
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
	msrv.Register(serverd.Service{
		Name:     "backend-builder.service",
		Start:    builder.Start,
		Shutdown: func() { builder.Shutdown() },
		Ping:     func() error { return builder.PingAll() },
	})
	return builder, nil
}

func createServices(finder archive.BackendFinder, msrv *serverd.Manager, logger yalogi.Logger) (*service.Builder, error) {
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
	msrv.Register(serverd.Service{
		Name:     "service-builder.service",
		Start:    builder.Start,
		Shutdown: func() { builder.Shutdown() },
	})
	return builder, nil
}

func createArchiveEventAPI(gsrv *grpc.Server, finder archive.ServiceFinder, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("api-archive").(*iconfig.ArchiveAPICfg)
	if cfgArchive.Event != "" {
		logger.Infof("creating and registering archive event api service")
		gsvc, err := ifactory.ArchiveEventAPI(cfgArchive, finder)
		if err != nil {
			return err
		}
		eventapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}

func createArchiveDNSAPI(gsrv *grpc.Server, finder archive.ServiceFinder, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("api-archive").(*iconfig.ArchiveAPICfg)
	if cfgArchive.DNS != "" {
		logger.Infof("creating and registering archive dns api service")
		gsvc, err := ifactory.ArchiveDNSAPI(cfgArchive, finder)
		if err != nil {
			return err
		}
		dnsapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}

func createArchiveTLSAPI(gsrv *grpc.Server, finder archive.ServiceFinder, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("api-archive").(*iconfig.ArchiveAPICfg)
	if cfgArchive.TLS != "" {
		logger.Infof("creating and registering archive tls api service")
		gsvc, err := ifactory.ArchiveTLSAPI(cfgArchive, finder)
		if err != nil {
			return err
		}
		tlsapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}
