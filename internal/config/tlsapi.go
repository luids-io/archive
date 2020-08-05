// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// ArchiveTLSAPICfg stores archive service preferences
type ArchiveTLSAPICfg struct {
	Enable  bool
	Log     bool
	Service string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ArchiveTLSAPICfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.BoolVar(&cfg.Enable, aprefix+"enable", cfg.Enable, "Enable tls archive api.")
	pflag.BoolVar(&cfg.Log, aprefix+"log", cfg.Log, "Enable log in service.")
	pflag.StringVar(&cfg.Service, aprefix+"service", cfg.Service, "Service id archive tls.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ArchiveTLSAPICfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"enable")
	util.BindViper(v, aprefix+"log")
	util.BindViper(v, aprefix+"service")
}

// FromViper fill values from viper
func (cfg *ArchiveTLSAPICfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Enable = v.GetBool(aprefix + "enable")
	cfg.Log = v.GetBool(aprefix + "log")
	cfg.Service = v.GetString(aprefix + "service")
}

// Empty returns true if configuration is empty
func (cfg ArchiveTLSAPICfg) Empty() bool {
	return !cfg.Enable
}

// Validate checks that configuration is ok
func (cfg ArchiveTLSAPICfg) Validate() error {
	if cfg.Service == "" {
		return fmt.Errorf("service must be defined")
	}
	return nil
}

// Dump configuration
func (cfg ArchiveTLSAPICfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
