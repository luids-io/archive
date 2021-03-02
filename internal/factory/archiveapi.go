// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"errors"
	"fmt"

	"github.com/luids-io/api/dnsutil"
	dnsapi "github.com/luids-io/api/dnsutil/grpc/archive"
	"github.com/luids-io/api/event"
	eventapi "github.com/luids-io/api/event/grpc/archive"
	"github.com/luids-io/api/tlsutil"
	tlsapi "github.com/luids-io/api/tlsutil/grpc/archive"
	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/core/yalogi"
)

// ArchiveEventAPI creates grpc service
func ArchiveEventAPI(cfg *config.ArchiveEventAPICfg, finder *archive.Builder, logger yalogi.Logger) (*eventapi.Service, error) {
	if !cfg.Enable {
		return nil, errors.New("event api disabled")
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getService(cfg.Service, archive.EventAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'eventapi' service: %v", err)
	}
	c, ok := svc.(event.Archiver)
	if !ok {
		return nil, fmt.Errorf("can't cast id '%s' to event.Archiver", cfg.Service)
	}
	if !cfg.Log {
		logger = yalogi.LogNull
	}
	return eventapi.NewService(c, eventapi.SetServiceLogger(logger)), nil
}

// ArchiveDNSAPI creates grpc service
func ArchiveDNSAPI(cfg *config.ArchiveDNSAPICfg, finder *archive.Builder, logger yalogi.Logger) (*dnsapi.Service, error) {
	if !cfg.Enable {
		return nil, errors.New("event api disabled")
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getService(cfg.Service, archive.DNSAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'dnsapi' service: %v", err)
	}
	c, ok := svc.(dnsutil.Archiver)
	if !ok {
		return nil, fmt.Errorf("can't cast id '%s' to dnsutil.Archiver", cfg.Service)
	}
	if !cfg.Log {
		logger = yalogi.LogNull
	}
	return dnsapi.NewService(c, dnsapi.SetServiceLogger(logger)), nil
}

// ArchiveTLSAPI creates grpc service
func ArchiveTLSAPI(cfg *config.ArchiveTLSAPICfg, finder *archive.Builder, logger yalogi.Logger) (*tlsapi.Service, error) {
	if !cfg.Enable {
		return nil, errors.New("event api disabled")
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getService(cfg.Service, archive.TLSAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'tlsapi' service: %v", err)
	}
	c, ok := svc.(tlsutil.Archiver)
	if !ok {
		return nil, fmt.Errorf("can't cast id '%s' to tlsutil.Archiver", cfg.Service)
	}
	if !cfg.Log {
		logger = yalogi.LogNull
	}
	return tlsapi.NewService(c, tlsapi.SetServiceLogger(logger)), nil
}
