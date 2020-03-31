// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// ArchiveAPICfg stores archive service preferences
type ArchiveAPICfg struct {
	Event string
	DNS   string
	TLS   string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ArchiveAPICfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.Event, aprefix+"event", cfg.Event, "Service id archive events.")
	pflag.StringVar(&cfg.DNS, aprefix+"dns", cfg.DNS, "Service id archive dns.")
	pflag.StringVar(&cfg.TLS, aprefix+"tls", cfg.TLS, "Service id archive tls.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ArchiveAPICfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"event")
	util.BindViper(v, aprefix+"dns")
	util.BindViper(v, aprefix+"tls")
}

// FromViper fill values from viper
func (cfg *ArchiveAPICfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.Event = v.GetString(aprefix + "event")
	cfg.DNS = v.GetString(aprefix + "dns")
	cfg.TLS = v.GetString(aprefix + "tls")
}

// Empty returns true if configuration is empty
func (cfg ArchiveAPICfg) Empty() bool {
	if cfg.Event != "" || cfg.DNS != "" || cfg.TLS != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg ArchiveAPICfg) Validate() error {
	if cfg.Event == "" && cfg.DNS == "" && cfg.TLS == "" {
		return fmt.Errorf("services must be defined")
	}
	return nil
}

// Dump configuration
func (cfg ArchiveAPICfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
