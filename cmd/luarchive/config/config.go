// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"

	iconfig "github.com/luids-io/archive/internal/config"
	cconfig "github.com/luids-io/common/config"
	"github.com/luids-io/core/goconfig"
)

// Default returns the default configuration
func Default(program string) *goconfig.Config {
	cfg, err := goconfig.New(program,
		goconfig.Section{
			Name:     "archive",
			Required: true,
			Short:    false,
			Data:     &iconfig.ArchiverCfg{},
		},
		goconfig.Section{
			Name:     "service.event.archive",
			Required: false,
			Short:    false,
			Data:     &iconfig.ArchiveEventAPICfg{Log: true},
		},
		goconfig.Section{
			Name:     "service.dnsutil.archive",
			Required: false,
			Short:    false,
			Data:     &iconfig.ArchiveDNSAPICfg{Log: true},
		},
		goconfig.Section{
			Name:     "service.tlsutil.archive",
			Required: false,
			Short:    false,
			Data:     &iconfig.ArchiveTLSAPICfg{Log: true},
		},
		goconfig.Section{
			Name:     "server",
			Required: true,
			Short:    true,
			Data: &cconfig.ServerCfg{
				ListenURI: "tcp://127.0.0.1:5821",
			},
		},
		goconfig.Section{
			Name:     "log",
			Required: true,
			Data: &cconfig.LoggerCfg{
				Level: "info",
			},
		},
		goconfig.Section{
			Name:     "health",
			Required: false,
			Data:     &cconfig.HealthCfg{},
		},
	)
	if err != nil {
		panic(err)
	}
	// add aditional validators
	cfg.AddValidator(func(cfg *goconfig.Config) error {
		noEvent := cfg.Data("service.event.archive").Empty()
		noDNS := cfg.Data("service.dnsutil.archive").Empty()
		noTLS := cfg.Data("service.tlsutil.archive").Empty()
		if noEvent && noDNS && noTLS {
			return errors.New("enable service is required")
		}
		return nil
	})
	return cfg
}
