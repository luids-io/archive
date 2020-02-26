// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"github.com/luisguillenc/goconfig"

	iconfig "github.com/luids-io/archive/internal/config"
	cconfig "github.com/luids-io/common/config"
)

// Default returns the default configuration
func Default(program string) *goconfig.Config {
	cfg, err := goconfig.New(program,
		goconfig.Section{
			Name:     "service",
			Required: true,
			Short:    false,
			Data:     &iconfig.ServiceCfg{},
		},
		goconfig.Section{
			Name:     "backend",
			Required: false,
			Short:    false,
			Data:     &iconfig.BackendCfg{},
		},
		goconfig.Section{
			Name:     "grpc-archive",
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
	return cfg
}
