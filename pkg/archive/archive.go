// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

// Package archive provides a generic system to create archive api services.
//
// This package is a work in progress and makes no API stability promises.
package archive

// API stores archive APIs available
type API int

// List of valid archive APIs
const (
	EventAPI API = iota
	DNSAPI
	TLSAPI
)

// Backend interface is a container that stores backend information
type Backend interface {
	ID() string
	Class() string
	Session() interface{}
	Ping() error
}

// Service interface for archive services
type Service interface {
	ID() string
	Class() string
	Implements() []API
}
