// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package main

import (
	"fmt"
	"os"

	"github.com/luisguillenc/serverd"
	"github.com/spf13/pflag"

	"github.com/luids-io/archive/cmd/luarchive/config"
)

//Variables for version output
var (
	Program  = "luarchive"
	Build    = "unknown"
	Version  = "unknown"
	Revision = "unknown"
)

var (
	cfg        = config.Default(Program)
	configFile = ""
	version    = false
	help       = false
	debug      = false
	dryRun     = false
)

func init() {
	//config mapped params
	cfg.PFlags()
	//behaviour params
	pflag.StringVar(&configFile, "config", configFile, "Use explicit config file.")
	pflag.BoolVar(&version, "version", version, "Show version.")
	pflag.BoolVarP(&help, "help", "h", help, "Show this help.")
	pflag.BoolVar(&debug, "debug", debug, "Enable debug.")
	pflag.BoolVar(&dryRun, "dry-run", dryRun, "Check connections but don't start service.")
	pflag.Parse()
}

func main() {
	if version {
		fmt.Printf("version: %s\nrevision: %s\nbuild: %s\n", Version, Revision, Build)
		os.Exit(0)
	}
	if help {
		pflag.Usage()
		os.Exit(0)
	}

	// load configuration
	err := cfg.LoadIfFile(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//creates logger
	logger, err := createLogger(debug)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// echo version and config
	logger.Infof("%s (version: %s build: %s)", Program, Version, Build)
	if debug {
		logger.Debugf("configuration dump: %v", cfg)
	}

	// creates main server manager
	srv := serverd.New(serverd.SetLogger(logger))

	// create backends and archive services
	backends, err := createBackends(srv, logger)
	if err != nil {
		logger.Fatalf("couldn't create backends: %v", err)
	}
	services, err := createServices(srv, backends, logger)
	if err != nil {
		logger.Fatalf("couldn't create backends: %v", err)
	}

	if dryRun {
		fmt.Println("configuration seems ok")
		os.Exit(0)
	}

	// create grpc server
	gsrv, err := createArchiverSrv(srv, logger)
	if err != nil {
		logger.Fatalf("couldn't create grpc server: %v", err)
	}

	// create grpc services and register in server
	err = registerServices(srv, gsrv, services, logger)
	if err != nil {
		logger.Fatalf("couldn't create grpc services in server: %v", err)
	}

	// creates health server
	err = createHealthSrv(srv, logger)
	if err != nil {
		logger.Fatalf("couldn't create health server: %v", err)
	}

	//run server
	err = srv.Run()
	if err != nil {
		logger.Errorf("running server: %v", err)
	}
	logger.Infof("%s finished", Program)
}
