// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luids-io/archive/pkg/archive"
)

// BuildFn defines a function that constructs a Backend
type BuildFn func(b *Builder, def Definition) (archive.Backend, error)

// Definition stores configuration definition of backends
type Definition struct {
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
	// TLS defines the configuration of client protocol
	TLS *ConfigTLS `json:"tls,omitempty"`
	// Opts custom options
	Opts map[string]interface{} `json:"opts,omitempty"`
}

//ConfigTLS stores information used in TLS connections
type ConfigTLS struct {
	CertFile     string `json:"certfile,omitempty"`
	KeyFile      string `json:"keyfile,omitempty"`
	ServerName   string `json:"servername,omitempty"`
	ServerCert   string `json:"servercert,omitempty"`
	CACert       string `json:"cacert,omitempty"`
	UseSystemCAs bool   `json:"systemca"`
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
