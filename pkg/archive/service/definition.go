// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luids-io/archive/pkg/archive"
)

// BuildFn defines a function that constructs a Service
type BuildFn func(b *Builder, def Definition) (*archive.Service, error)

// Definition stores configuration definition of services
type Definition struct {
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

// DefsFromFile creates a slice from a file in json format.
func DefsFromFile(path string) ([]Definition, error) {
	var defs []Definition
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

// RegisterBuilder registers a builder
func RegisterBuilder(class string, builder BuildFn) {
	registryBuilder[class] = builder
}

var registryBuilder map[string]BuildFn

func init() {
	registryBuilder = make(map[string]BuildFn)
}
