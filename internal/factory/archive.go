// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"errors"
	"fmt"

	"github.com/luisguillenc/yalogi"

	dnsapi "github.com/luids-io/api/dnsutil/archive"
	eventapi "github.com/luids-io/api/event/archive"
	tlsapi "github.com/luids-io/api/tlsutil/archive"
	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/core/dnsutil"
	"github.com/luids-io/core/event"
	"github.com/luids-io/core/tlsutil"
)

// EventAPIService creates grpc service
func EventAPIService(cfg *config.ArchiveCfg, finder archive.ServiceFinder, logger yalogi.Logger) (*eventapi.Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getArchiveService(cfg.EventAPI, archive.EventAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'eventapi' service: %v", err)
	}
	c, ok := svc.(event.Archiver)
	if !ok {
		logger.Fatalf("can't cast id '%s' to event.Archiver", cfg.EventAPI)
	}
	return eventapi.NewService(c), nil
}

// DNSAPIService creates grpc service
func DNSAPIService(cfg *config.ArchiveCfg, finder archive.ServiceFinder, logger yalogi.Logger) (*dnsapi.Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getArchiveService(cfg.DNSAPI, archive.DNSAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'dnsapi' service: %v", err)
	}
	c, ok := svc.(dnsutil.Archiver)
	if !ok {
		logger.Fatalf("can't cast id '%s' to dnsutil.Archiver", cfg.DNSAPI)
	}
	return dnsapi.NewService(c), nil
}

// TLSAPIService creates grpc service
func TLSAPIService(cfg *config.ArchiveCfg, finder archive.ServiceFinder, logger yalogi.Logger) (*tlsapi.Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getArchiveService(cfg.TLSAPI, archive.TLSAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'tlsapi' service: %v", err)
	}
	c, ok := svc.(tlsutil.Archiver)
	if !ok {
		logger.Fatalf("can't cast id '%s' to tlsutil.Archiver", cfg.TLSAPI)
	}
	return tlsapi.NewService(c), nil
}

func getArchiveService(name string, api archive.API, finder archive.ServiceFinder) (archive.Service, error) {
	if name == "" {
		return nil, errors.New("service id is empty")
	}
	svc, ok := finder.FindServiceByID(name)
	if !ok {
		return nil, fmt.Errorf("can't find service with id '%s'", name)
	}
	if !implements(svc, api) {
		return nil, fmt.Errorf("service '%s' don't implements api", name)
	}
	return svc, nil
}

func implements(svc archive.Service, api archive.API) bool {
	for _, v := range svc.Implements() {
		if v == api {
			return true
		}
	}
	return false
}
