// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package mongodb

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo"

	"github.com/luids-io/archive/pkg/archive"
)

// Builder returns a builder function.
func Builder() archive.BuildBackendFn {
	return func(b *archive.Builder, def archive.BackendDef) (archive.Backend, error) {
		if def.URL == "" {
			return nil, errors.New("'url' is required")
		}
		// create session with mgo
		session, err := mgo.Dial(def.URL)
		if err != nil {
			return nil, fmt.Errorf("dialing with mongodb '%s': %v", def.URL, err)
		}
		backend := &mdbBackend{id: def.ID, session: session}
		b.OnShutdown(func() error {
			session.Close()
			return nil
		})
		return backend, nil
	}
}

func init() {
	archive.RegisterBackendBuilder(BackendClass, Builder())
}
