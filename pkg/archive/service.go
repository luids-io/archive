// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package archive

// API stores archive APIs available
type API int

// List of valid archive APIs
const (
	EventAPI API = iota
	DNSAPI
	TLSAPI
)

// Service container stores service information
type Service struct {
	ID     string
	Class  string
	API    API
	Object interface{}
}

// ServiceFinder interface for services
type ServiceFinder interface {
	FindServiceByID(string) (*Service, bool)
	FindAllServices() []*Service
}
