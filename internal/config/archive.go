// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"fmt"

	"github.com/luids-io/common/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ArchiveCfg stores archive service preferences
type ArchiveCfg struct {
	EventAPI string
	DNSAPI   string
	TLSAPI   string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ArchiveCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.EventAPI, aprefix+"eventapi", cfg.EventAPI, "Service id archive events api.")
	pflag.StringVar(&cfg.DNSAPI, aprefix+"dnsapi", cfg.DNSAPI, "Service id archive dns api.")
	pflag.StringVar(&cfg.TLSAPI, aprefix+"tlsapi", cfg.TLSAPI, "Service id archive tls api.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ArchiveCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"eventapi")
	util.BindViper(v, aprefix+"dnsapi")
	util.BindViper(v, aprefix+"tlsapi")
}

// FromViper fill values from viper
func (cfg *ArchiveCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.EventAPI = v.GetString(aprefix + "eventapi")
	cfg.DNSAPI = v.GetString(aprefix + "dnsapi")
	cfg.TLSAPI = v.GetString(aprefix + "tlsapi")
}

// Empty returns true if configuration is empty
func (cfg ArchiveCfg) Empty() bool {
	if cfg.EventAPI != "" || cfg.DNSAPI != "" || cfg.TLSAPI != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg ArchiveCfg) Validate() error {
	if cfg.EventAPI == "" && cfg.DNSAPI == "" && cfg.TLSAPI == "" {
		return fmt.Errorf("api must be defined")
	}
	return nil
}

// Dump configuration
func (cfg ArchiveCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
