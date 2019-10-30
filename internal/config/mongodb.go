// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// MongoDBCfg stores mongodb storage settings
type MongoDBCfg struct {
	URL      string
	Database string
	Prefix   string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *MongoDBCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	pflag.StringVar(&cfg.URL, aprefix+"url", cfg.URL, "URL server.")
	pflag.StringVar(&cfg.Database, aprefix+"db", cfg.Database, "Database for mongodb.")
	pflag.StringVar(&cfg.Prefix, aprefix+"prefix", cfg.Prefix, "Collections prefix.")
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *MongoDBCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"url")
	util.BindViper(v, aprefix+"db")
	util.BindViper(v, aprefix+"prefix")
}

// FromViper fill values from viper
func (cfg *MongoDBCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.URL = v.GetString(aprefix + "url")
	cfg.Database = v.GetString(aprefix + "db")
	cfg.Prefix = v.GetString(aprefix + "prefix")
}

// Empty returns true if configuration is empty
func (cfg MongoDBCfg) Empty() bool {
	if cfg.URL != "" {
		return false
	}
	if cfg.Database != "" {
		return false
	}
	if cfg.Prefix != "" {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg MongoDBCfg) Validate() error {
	if cfg.URL == "" {
		return errors.New("url is required")
	}
	_, err := mgo.ParseURL(cfg.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	return nil
}

// Dump configuration
func (cfg MongoDBCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
