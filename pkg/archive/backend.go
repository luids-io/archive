// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package archive

// Backend container stores backend information
type Backend interface {
	GetClass() string
	GetSession() interface{}
	Ping() error
}

// BackendFinder interface for backends
type BackendFinder interface {
	FindBackendByID(string) (Backend, bool)
}
