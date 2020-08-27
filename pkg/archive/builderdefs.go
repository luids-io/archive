// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package archive

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luids-io/core/grpctls"
)

// BuildBackendFn defines a function that constructs a Backend
type BuildBackendFn func(b *Builder, def BackendDef) (Backend, error)

// BuildServiceFn defines a function that constructs a Service
type BuildServiceFn func(b *Builder, def ServiceDef) (Service, error)

// BackendDef stores configuration definition of backends
type BackendDef struct {
	// ID must exist and be unique
	ID string `json:"id"`
	// Class stores the driver
	Class string `json:"class"`
	// Disabled flag
	Disabled bool `json:"disabled"`
	// Name or description
	Name string `json:"name,omitempty"`
	// URL provides the url to backend
	URL string `json:"url,omitempty"`
	// Client configuration
	Client *grpctls.ClientCfg `json:"tls,omitempty"`
	// Opts custom options
	Opts map[string]interface{} `json:"opts,omitempty"`
}

// ClientCfg returns a copy of client configuration.
// It returns an empty struct if a null pointer is stored.
func (def BackendDef) ClientCfg() grpctls.ClientCfg {
	if def.Client == nil {
		return grpctls.ClientCfg{}
	}
	return *def.Client
}

// ServiceDef stores configuration definition of services
type ServiceDef struct {
	// ID must exist and be unique
	ID string `json:"id"`
	// Class stores the services class
	Class string `json:"class"`
	// Disabled flag
	Disabled bool `json:"disabled"`
	// Name or description
	Name string `json:"name,omitempty"`
	// Backend ID
	Backend string `json:"backend,omitempty"`
	// Opts custom options
	Opts map[string]interface{} `json:"opts,omitempty"`
}

// BackendDefsFromFile creates a slice from a file in json format.
func BackendDefsFromFile(path string) ([]BackendDef, error) {
	var defs []BackendDef
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("opening file '%s': %v", path, err)
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading file '%s': %v", path, err)
	}
	err = json.Unmarshal(byteValue, &defs)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling from json file '%s': %v", path, err)
	}
	return defs, nil
}

// ServiceDefsFromFile creates a slice from a file in json format.
func ServiceDefsFromFile(path string) ([]ServiceDef, error) {
	var defs []ServiceDef
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("opening file '%s': %v", path, err)
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading file '%s': %v", path, err)
	}
	err = json.Unmarshal(byteValue, &defs)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling from json file '%s': %v", path, err)
	}
	return defs, nil
}

// RegisterBackendBuilder registers a builder
func RegisterBackendBuilder(class string, builder BuildBackendFn) {
	regBackendBuilder[class] = builder
}

// RegisterServiceBuilder registers a builder
func RegisterServiceBuilder(class string, builder BuildServiceFn) {
	regServiceBuilder[class] = builder
}

var regBackendBuilder map[string]BuildBackendFn
var regServiceBuilder map[string]BuildServiceFn

func init() {
	regBackendBuilder = make(map[string]BuildBackendFn)
	regServiceBuilder = make(map[string]BuildServiceFn)
}
