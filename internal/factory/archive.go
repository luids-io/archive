// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/archive/builder"
	"github.com/luids-io/common/util"
	"github.com/luids-io/core/utils/yalogi"
)

// ArchiveBuilder is a factory
func ArchiveBuilder(cfg *config.ArchiverCfg, logger yalogi.Logger) (*builder.Builder, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	b := builder.NewBuilder(builder.SetLogger(logger))
	return b, nil
}

//Backends creates backends from configuration files
func Backends(cfg *config.ArchiverCfg, b *builder.Builder, logger yalogi.Logger) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("bad config: %v", err)
	}
	dbfiles, err := util.GetFilesDB("json", cfg.BackendFiles, cfg.BackendDirs)
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
		_, err := b.BuildBackend(def)
		if err != nil {
			return fmt.Errorf("creating '%s': %v", def.ID, err)
		}
	}
	return nil
}

func loadBackendDefs(dbFiles []string) ([]builder.BackendDef, error) {
	loadedDB := make([]builder.BackendDef, 0)
	for _, file := range dbFiles {
		entries, err := builder.BackendDefsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("couln't load database: %v", err)
		}
		loadedDB = append(loadedDB, entries...)
	}
	return loadedDB, nil
}

//Services creates services from configuration files
func Services(cfg *config.ArchiverCfg, b *builder.Builder, logger yalogi.Logger) error {
	err := cfg.Validate()
	if err != nil {
		return fmt.Errorf("bad config: %v", err)
	}
	dbfiles, err := util.GetFilesDB("json", cfg.ServiceFiles, cfg.ServiceDirs)
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
		_, err := b.BuildService(def)
		if err != nil {
			return fmt.Errorf("creating '%s': %v", def.ID, err)
		}
	}
	return nil
}

func loadServiceDefs(dbFiles []string) ([]builder.ServiceDef, error) {
	loadedDB := make([]builder.ServiceDef, 0)
	for _, file := range dbFiles {
		entries, err := builder.ServiceDefsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("couln't load database: %v", err)
		}
		loadedDB = append(loadedDB, entries...)
	}
	return loadedDB, nil
}
