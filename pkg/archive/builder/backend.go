// Copyright 2020 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luids-io/archive/pkg/archive"
)

// BuildBackendFn defines a function that constructs a Backend
type BuildBackendFn func(b *Builder, def BackendDef) (archive.Backend, error)

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

// RegisterBackendBuilder registers a builder
func RegisterBackendBuilder(class string, builder BuildBackendFn) {
	regBackendBuilder[class] = builder
}

var regBackendBuilder map[string]BuildBackendFn

func init() {
	regBackendBuilder = make(map[string]BuildBackendFn)
}
