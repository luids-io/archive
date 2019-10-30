// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// ArchiverCfg stores generic archive settings
type ArchiverCfg struct {
	Backend string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ArchiverCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	if short {
		pflag.StringVarP(&cfg.Backend, aprefix+"backend", "b", cfg.Backend, "Storage backend.")
	} else {
		pflag.StringVar(&cfg.Backend, aprefix+"backend", cfg.Backend, "Storage backend.")
	}
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ArchiverCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"backend")
}

// FromViper fill values from viper
func (cfg *ArchiverCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Backend = v.GetString(aprefix + "backend")
}

// Empty returns true if configuration is empty
func (cfg ArchiverCfg) Empty() bool {
	if cfg.Backend != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg ArchiverCfg) Validate() error {
	if cfg.Backend == "" {
		return errors.New("backend is required")
	}
	return nil
}

// Dump configuration
func (cfg ArchiverCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
