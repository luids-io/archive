// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	// backends
	_ "github.com/luids-io/archive/pkg/archive/backends/mongodb"

	// services
	_ "github.com/luids-io/archive/pkg/archive/services/dnsmdb"
	_ "github.com/luids-io/archive/pkg/archive/services/eventmdb"
	_ "github.com/luids-io/archive/pkg/archive/services/tlsmdb"
)
