// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package archive

// Backend container stores backend information
type Backend struct {
	ID     string
	Class  string
	Object interface{}
}

// BackendFinder interface for backends
type BackendFinder interface {
	Backend(string) (*Backend, bool)
	Backends() []*Backend
}
