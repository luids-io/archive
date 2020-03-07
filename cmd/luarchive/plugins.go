// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	// backends
	_ "github.com/luids-io/archive/pkg/archive/backend/mongodb"

	// services
	_ "github.com/luids-io/archive/pkg/archive/service/dnsmdb"
	_ "github.com/luids-io/archive/pkg/archive/service/eventmdb"
	_ "github.com/luids-io/archive/pkg/archive/service/tlsmdb"
)
