// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"github.com/luisguillenc/goconfig"

	cconfig "github.com/luids-io/common/config"
	iconfig "github.com/luids-io/archive/internal/config"
)

// Default returns the default configuration
func Default(program string) *goconfig.Config {
	cfg, err := goconfig.New(program,
		goconfig.Section{
			Name:     "",
			Required: true,
			Short:    true,
			Data:     &iconfig.ArchiverCfg{Backend: "mongodb"},
		},
		goconfig.Section{
			Name:     "mongodb",
			Required: false,
			Data: &iconfig.MongoDBCfg{
				Database: "dnsdb",
				URL:      "127.0.0.1:27017",
			},
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
