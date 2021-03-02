// Copyright 2021 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"errors"
	"fmt"

	"github.com/luids-io/api/dnsutil"
	dnsapi "github.com/luids-io/api/dnsutil/grpc/finder"
	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/core/yalogi"
)

// FinderDNSAPI creates grpc service
func FinderDNSAPI(cfg *config.FinderDNSAPICfg, finder *archive.Builder, logger yalogi.Logger) (*dnsapi.Service, error) {
	if !cfg.Enable {
		return nil, errors.New("dns finder api disabled")
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	svc, err := getService(cfg.Service, archive.DNSAPI, finder)
	if err != nil {
		return nil, fmt.Errorf("'dnsapi' service: %v", err)
	}
	f, ok := svc.(dnsutil.Finder)
	if !ok {
		return nil, fmt.Errorf("can't cast id '%s' to dnsutil.Finder", cfg.Service)
	}
	if !cfg.Log {
		logger = yalogi.LogNull
	}
	return dnsapi.NewService(f, dnsapi.SetServiceLogger(logger)), nil
}
