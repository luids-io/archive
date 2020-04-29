// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"

	iconfig "github.com/luids-io/archive/internal/config"
	cconfig "github.com/luids-io/common/config"
	"github.com/luids-io/core/utils/goconfig"
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
			Name:     "archive.api.event",
			Required: false,
			Short:    false,
			Data:     &iconfig.ArchiveEventAPICfg{},
		},
		goconfig.Section{
			Name:     "archive.api.dns",
			Required: false,
			Short:    false,
			Data:     &iconfig.ArchiveDNSAPICfg{},
		},
		goconfig.Section{
			Name:     "archive.api.tls",
			Required: false,
			Short:    false,
			Data:     &iconfig.ArchiveTLSAPICfg{},
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
		noEvent := cfg.Data("archive.api.event").Empty()
		noDNS := cfg.Data("archive.api.dns").Empty()
		noTLS := cfg.Data("archive.api.tls").Empty()
		if noEvent && noDNS && noTLS {
			return errors.New("'archive.api' section is required")
		}
		return nil
	})
	return cfg
}
