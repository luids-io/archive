// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/luids-io/common/util"
)

// ArchiverCfg stores archiver preferences
type ArchiverCfg struct {
	BackendDirs  []string
	BackendFiles []string
	ServiceDirs  []string
	ServiceFiles []string
}

// SetPFlags setups posix flags for commandline configuration
func (cfg *ArchiverCfg) SetPFlags(short bool, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	if short {
		pflag.StringSliceVarP(&cfg.BackendDirs, aprefix+"backend.dirs", "B", cfg.BackendDirs, "Backend dirs.")
		pflag.StringSliceVarP(&cfg.BackendFiles, aprefix+"backend.files", "b", cfg.BackendFiles, "Backend files.")
		pflag.StringSliceVarP(&cfg.ServiceDirs, aprefix+"service.dirs", "S", cfg.ServiceDirs, "Service dirs.")
		pflag.StringSliceVarP(&cfg.ServiceFiles, aprefix+"service.files", "s", cfg.ServiceFiles, "Service files.")
	} else {
		pflag.StringSliceVar(&cfg.BackendDirs, aprefix+"backend.dirs", cfg.BackendDirs, "Backend dirs.")
		pflag.StringSliceVar(&cfg.BackendFiles, aprefix+"backend.files", cfg.BackendFiles, "Backend files.")
		pflag.StringSliceVar(&cfg.ServiceDirs, aprefix+"service.dirs", cfg.ServiceDirs, "Service dirs.")
		pflag.StringSliceVar(&cfg.ServiceFiles, aprefix+"service.files", cfg.ServiceFiles, "Service files.")
	}
}

// BindViper setups posix flags for commandline configuration and bind to viper
func (cfg *ArchiverCfg) BindViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	util.BindViper(v, aprefix+"backend.dirs")
	util.BindViper(v, aprefix+"backend.files")
	util.BindViper(v, aprefix+"service.dirs")
	util.BindViper(v, aprefix+"service.files")
}

// FromViper fill values from viper
func (cfg *ArchiverCfg) FromViper(v *viper.Viper, prefix string) {
	aprefix := ""
	if prefix != "" {
		aprefix = prefix + "."
	}
	cfg.BackendDirs = v.GetStringSlice(aprefix + "backend.dirs")
	cfg.BackendFiles = v.GetStringSlice(aprefix + "backend.files")
	cfg.ServiceDirs = v.GetStringSlice(aprefix + "service.dirs")
	cfg.ServiceFiles = v.GetStringSlice(aprefix + "service.files")
}

// Empty returns true if configuration is empty
func (cfg ArchiverCfg) Empty() bool {
	if len(cfg.BackendDirs) > 0 {
		return false
	}
	if len(cfg.BackendFiles) > 0 {
		return false
	}
	if len(cfg.ServiceDirs) > 0 {
		return false
	}
	if len(cfg.ServiceFiles) > 0 {
		return false
	}
	return true
}

// Validate checks that configuration is ok
func (cfg ArchiverCfg) Validate() error {
	// parse config files
	if len(cfg.BackendFiles) == 0 && len(cfg.BackendDirs) == 0 {
		return errors.New("config backend required")
	}
	if len(cfg.ServiceFiles) == 0 && len(cfg.ServiceDirs) == 0 {
		return errors.New("config service required")
	}
	for _, file := range cfg.BackendFiles {
		if !strings.HasSuffix(file, ".json") {
			return fmt.Errorf("config file '%s' without .json extension", file)
		}
		if !util.FileExists(file) {
			return fmt.Errorf("config file '%v' doesn't exists", file)
		}
	}
	for _, dir := range cfg.BackendDirs {
		if !util.DirExists(dir) {
			return fmt.Errorf("config dir '%v' doesn't exists", dir)
		}
	}
	for _, file := range cfg.ServiceFiles {
		if !strings.HasSuffix(file, ".json") {
			return fmt.Errorf("config file '%s' without .json extension", file)
		}
		if !util.FileExists(file) {
			return fmt.Errorf("config file '%v' doesn't exists", file)
		}
	}
	for _, dir := range cfg.ServiceDirs {
		if !util.DirExists(dir) {
			return fmt.Errorf("config dir '%v' doesn't exists", dir)
		}
	}
	return nil
}

// Dump configuration
func (cfg ArchiverCfg) Dump() string {
	return fmt.Sprintf("%+v", cfg)
}
