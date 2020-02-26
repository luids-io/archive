// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/archive/backend"
	"github.com/luids-io/common/util"
)

// BackendBuilder is a factory
func BackendBuilder(cfg *config.BackendCfg, logger yalogi.Logger) (*backend.Builder, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	b := backend.New(backend.SetLogger(logger))
	return b, nil
}

//Backends creates backends from configuration files
func Backends(cfg *config.BackendCfg, builder *backend.Builder, logger yalogi.Logger) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("bad config: %v", err)
	}
	dbfiles, err := util.GetFilesDB("json", cfg.ConfigFiles, cfg.ConfigDirs)
	if err != nil {
		return fmt.Errorf("loading dbfiles: %v", err)
	}
	defs, err := loadBackendDefs(dbfiles)
	if err != nil {
		return fmt.Errorf("loading dbfiles: %v", err)
	}
	for _, def := range defs {
		if def.Disabled {
			continue
		}
		_, err := builder.Build(def)
		if err != nil {
			return fmt.Errorf("creating '%s': %v", def.ID, err)
		}
	}
	return nil
}

func loadBackendDefs(dbFiles []string) ([]backend.Definition, error) {
	loadedDB := make([]backend.Definition, 0)
	for _, file := range dbFiles {
		entries, err := backend.DefsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("couln't load database: %v", err)
		}
		loadedDB = append(loadedDB, entries...)
	}
	return loadedDB, nil
}
