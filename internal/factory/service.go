// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/archive"
	"github.com/luids-io/archive/pkg/archive/service"
	"github.com/luids-io/common/util"
)

// ServiceBuilder is a factory
func ServiceBuilder(cfg *config.ServiceCfg, finder archive.BackendFinder, logger yalogi.Logger) (*service.Builder, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	b := service.NewBuilder(finder, service.SetLogger(logger))
	return b, nil
}

//Services creates services from configuration files
func Services(cfg *config.ServiceCfg, builder *service.Builder, logger yalogi.Logger) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("bad config: %v", err)
	}
	dbfiles, err := util.GetFilesDB("json", cfg.ConfigFiles, cfg.ConfigDirs)
	if err != nil {
		return fmt.Errorf("loading dbfiles: %v", err)
	}
	defs, err := loadServiceDefs(dbfiles)
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

func loadServiceDefs(dbFiles []string) ([]service.Definition, error) {
	loadedDB := make([]service.Definition, 0)
	for _, file := range dbFiles {
		entries, err := service.DefsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("couln't load database: %v", err)
		}
		loadedDB = append(loadedDB, entries...)
	}
	return loadedDB, nil
}
