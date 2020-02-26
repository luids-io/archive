// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/luids-io/common/util"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BackendCfg stores backend preferences
type BackendCfg struct {
	ConfigDirs  []string
	ConfigFiles []string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *BackendCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	if short {
		pflag.StringSliceVarP(&cfg.ConfigDirs, aprefix+"dirs", "B", cfg.ConfigDirs, "Backend dirs.")
		pflag.StringSliceVarP(&cfg.ConfigFiles, aprefix+"files", "b", cfg.ConfigFiles, "Backend files.")
	} else {
		pflag.StringSliceVar(&cfg.ConfigDirs, aprefix+"dirs", cfg.ConfigDirs, "Backend dirs.")
		pflag.StringSliceVar(&cfg.ConfigFiles, aprefix+"files", cfg.ConfigFiles, "backend files.")
	}
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *BackendCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"dirs")
	util.BindViper(v, aprefix+"files")
}

// FromViper fill values from viper
func (cfg *BackendCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.ConfigDirs = v.GetStringSlice(aprefix + "dirs")
	cfg.ConfigFiles = v.GetStringSlice(aprefix + "files")
}

// Empty returns true if configuration is empty
func (cfg BackendCfg) Empty() bool {
	if len(cfg.ConfigDirs) > 0 {
		return false
	}
	if len(cfg.ConfigFiles) > 0 {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg BackendCfg) Validate() error {
	// parse config files
	if len(cfg.ConfigFiles) == 0 && len(cfg.ConfigDirs) == 0 {
		return errors.New("config required")
	}
	for _, file := range cfg.ConfigFiles {
		if !strings.HasSuffix(file, ".json") {
			return fmt.Errorf("config file '%s' without .json extension", file)
		}
		if !util.FileExists(file) {
			return fmt.Errorf("config file '%v' doesn't exists", file)
		}
	}
	for _, dir := range cfg.ConfigDirs {
		if !util.DirExists(dir) {
			return fmt.Errorf("config dir '%v' doesn't exists", dir)
		}
	}
	return nil
}

// Dump configuration
func (cfg BackendCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
