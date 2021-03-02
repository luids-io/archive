// Copyright 2021 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// FinderDNSAPICfg stores finder service preferences
type FinderDNSAPICfg struct {
	Enable  bool
	Log     bool
	Service string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *FinderDNSAPICfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.BoolVar(&cfg.Enable, aprefix+"enable", cfg.Enable, "Enable dns finder api.")
	pflag.BoolVar(&cfg.Log, aprefix+"log", cfg.Log, "Enable log in service.")
	pflag.StringVar(&cfg.Service, aprefix+"service", cfg.Service, "Service id finder dns.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *FinderDNSAPICfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"enable")
	util.BindViper(v, aprefix+"log")
	util.BindViper(v, aprefix+"service")
}

// FromViper fill values from viper
func (cfg *FinderDNSAPICfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Enable = v.GetBool(aprefix + "enable")
	cfg.Log = v.GetBool(aprefix + "log")
	cfg.Service = v.GetString(aprefix + "service")
}

// Empty returns true if configuration is empty
func (cfg FinderDNSAPICfg) Empty() bool {
	return !cfg.Enable
}

// Validate checks that configuration is ok
func (cfg FinderDNSAPICfg) Validate() error {
	if cfg.Service == "" {
		return fmt.Errorf("service must be defined")
	}
	return nil
}

// Dump configuration
func (cfg FinderDNSAPICfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
