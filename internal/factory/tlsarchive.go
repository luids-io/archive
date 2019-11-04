// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/archive/internal/config"
	"github.com/luids-io/archive/pkg/tlsarchive/mongodb"
)

// TLSArchiveMDB is a factory
func TLSArchiveMDB(cfg *config.MongoDBCfg, logger yalogi.Logger) (*mongodb.Archiver, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid mongodb: %v", err)
	}
	session, err := mgo.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("dialing with mongodb '%s': %v", cfg.URL, err)
	}
	archiver := mongodb.New(session, cfg.Database,
		mongodb.SetLogger(logger),
		mongodb.SetPrefix(cfg.Prefix))
	return archiver, nil
}
