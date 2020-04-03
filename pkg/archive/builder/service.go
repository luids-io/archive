// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luids-io/archive/pkg/archive"
)

// BuildServiceFn defines a function that constructs a Service
type BuildServiceFn func(b *Builder, def ServiceDef) (archive.Service, error)

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

// RegisterServiceBuilder registers a builder
func RegisterServiceBuilder(class string, builder BuildServiceFn) {
	regServiceBuilder[class] = builder
}

var regServiceBuilder map[string]BuildServiceFn

func init() {
	regServiceBuilder = make(map[string]BuildServiceFn)
}
