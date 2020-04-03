// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package archive

// Backend container stores backend information
type Backend interface {
	Class() string
	Session() interface{}
	Ping() error
}
