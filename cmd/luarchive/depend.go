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
			Name:     fmt.Sprintf("health.[%s]", cfgHealth.ListenURI),
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
		Name:     fmt.Sprintf("server.[%s]", cfgServer.ListenURI),
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
		Name:     "archive",
		Start:    builder.Start,
		Shutdown: func() { builder.Shutdown() },
		Ping:     func() error { return builder.PingAll() },
	})
	return builder, nil
}

func createArchiveEventAPI(gsrv *grpc.Server, finder *builder.Builder, msrv *serverd.Manager, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("service.event.archive").(*iconfig.ArchiveEventAPICfg)
	if cfgArchive.Enable {
		gsvc, err := ifactory.ArchiveEventAPI(cfgArchive, finder, logger)
		if err != nil {
			return err
		}
		eventapi.RegisterServer(gsrv, gsvc)
		msrv.Register(serverd.Service{Name: "service.event.archive"})
	}
	return nil
}

func createArchiveDNSAPI(gsrv *grpc.Server, finder *builder.Builder, msrv *serverd.Manager, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("service.dnsutil.archive").(*iconfig.ArchiveDNSAPICfg)
	if cfgArchive.Enable {
		gsvc, err := ifactory.ArchiveDNSAPI(cfgArchive, finder, logger)
		if err != nil {
			return err
		}
		dnsapi.RegisterServer(gsrv, gsvc)
		msrv.Register(serverd.Service{Name: "service.dnsutil.archive"})
	}
	return nil
}

func createArchiveTLSAPI(gsrv *grpc.Server, finder *builder.Builder, msrv *serverd.Manager, logger yalogi.Logger) error {
	cfgArchive := cfg.Data("service.tlsutil.archive").(*iconfig.ArchiveTLSAPICfg)
	if cfgArchive.Enable {
		gsvc, err := ifactory.ArchiveTLSAPI(cfgArchive, finder, logger)
		if err != nil {
			return err
		}
		tlsapi.RegisterServer(gsrv, gsvc)
		msrv.Register(serverd.Service{Name: "service.tlsutil.archive"})
	}
	return nil
}
