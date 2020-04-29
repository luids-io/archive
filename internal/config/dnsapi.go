// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// ArchiveDNSAPICfg stores archive service preferences
type ArchiveDNSAPICfg struct {
	Enable  bool
	Service string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ArchiveDNSAPICfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.BoolVar(&cfg.Enable, aprefix+"enable", cfg.Enable, "Enable dns archive api.")
	pflag.StringVar(&cfg.Service, aprefix+"service", cfg.Service, "Service id archive dns.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ArchiveDNSAPICfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"enable")
	util.BindViper(v, aprefix+"service")
}

// FromViper fill values from viper
func (cfg *ArchiveDNSAPICfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Enable = v.GetBool(aprefix + "enable")
	cfg.Service = v.GetString(aprefix + "service")
}

// Empty returns true if configuration is empty
func (cfg ArchiveDNSAPICfg) Empty() bool {
	if !cfg.Enable {
		return true
	}
	if cfg.Service != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg ArchiveDNSAPICfg) Validate() error {
	if !cfg.Enable {
		return nil
	}
	if cfg.Service == "" {
		return fmt.Errorf("service must be defined")
	}
	return nil
}

// Dump configuration
func (cfg ArchiveDNSAPICfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
