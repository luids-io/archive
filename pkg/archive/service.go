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

// Service interface for archive services
type Service interface {
	GetClass() string
	Implements() []API
}

// ServiceFinder interface for services
type ServiceFinder interface {
	FindServiceByID(string) (Service, bool)
}
