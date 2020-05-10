// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	"fmt"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	dnsapi "github.com/luids-io/api/dnsutil/grpc/archive"
	eventapi "github.com/luids-io/api/event/grpc/archive"
	tlsapi "github.com/luids-io/api/tlsutil/grpc/archive"
	iconfig "github.com/luids-io/archive/internal/config"
	ifactory "github.com/luids-io/archive/internal/factory"
	"github.com/luids-io/archive/pkg/archive/builder"
	cconfig "github.com/luids-io/common/config"
	cfactory "github.com/luids-io/common/factory"
	"github.com/luids-io/core/serverd"
	"github.com/luids-io/core/yalogi"
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

func createServer(msrv *serverd.Manager) (*grpc.Server, error) {
	cfgServer := cfg.Data("server").(*cconfig.ServerCfg)
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

func createArchivers(msrv *serverd.Manager, logger yalogi.Logger) (*builder.Builder, error) {
	cfgArchive := cfg.Data("archive").(*iconfig.ArchiverCfg)
	builder, err := ifactory.ArchiveBuilder(cfgArchive, logger)
	if err != nil {
		return nil, err
	}
	//create backends
	err = ifactory.Backends(cfgArchive, builder, logger)
	if err != nil {
		return nil, err
	}
	//create services
	err = ifactory.Services(cfgArchive, builder, logger)
	if err != nil {
		return nil, err
	}
	msrv.Register(serverd.Service{
		Name:     "archive-builder.service",
		Start:    builder.Start,
		Shutdown: func() { builder.Shutdown() },
		Ping:     func() error { return builder.PingAll() },
	})
	return builder, nil
}

func createArchiveEventAPI(gsrv *grpc.Server, finder *builder.Builder, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("archive.api.event").(*iconfig.ArchiveEventAPICfg)
	if cfgArchive.Enable {
		logger.Infof("registering archive event api service")
		gsvc, err := ifactory.ArchiveEventAPI(cfgArchive, finder)
		if err != nil {
			return err
		}
		eventapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}

func createArchiveDNSAPI(gsrv *grpc.Server, finder *builder.Builder, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("archive.api.dns").(*iconfig.ArchiveDNSAPICfg)
	if cfgArchive.Enable {
		logger.Infof("registering archive dns api service")
		gsvc, err := ifactory.ArchiveDNSAPI(cfgArchive, finder)
		if err != nil {
			return err
		}
		dnsapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}

func createArchiveTLSAPI(gsrv *grpc.Server, finder *builder.Builder, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("archive.api.tls").(*iconfig.ArchiveTLSAPICfg)
	if cfgArchive.Enable {
		logger.Infof("registering archive tls api service")
		gsvc, err := ifactory.ArchiveTLSAPI(cfgArchive, finder)
		if err != nil {
			return err
		}
		tlsapi.RegisterServer(gsrv, gsvc)
		return nil
	}
	return nil
}
