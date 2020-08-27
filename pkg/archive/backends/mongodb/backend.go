// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package mongodb

import (
	"github.com/globalsign/mgo"
)

// BackendClass registered.
const BackendClass = "mongodb"

//mdbBackend implements archive.Backend interface
type mdbBackend struct {
	id      string
	session *mgo.Session
}

func (b *mdbBackend) ID() string {
	return b.id
}

func (b *mdbBackend) Class() string {
	return BackendClass
}

func (b *mdbBackend) Session() interface{} {
	return b.session
}

func (b *mdbBackend) Ping() error {
	return b.session.Ping()
}
