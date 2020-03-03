// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package mongodb

import (
	"github.com/globalsign/mgo"
)

//mdbBackend implements archive.Backend interface
type mdbBackend struct {
	session *mgo.Session
}

func (b *mdbBackend) GetClass() string {
	return BackendClass
}

func (b *mdbBackend) GetSession() interface{} {
	return b.session
}

func (b *mdbBackend) Ping() error {
	return b.session.Ping()
}
