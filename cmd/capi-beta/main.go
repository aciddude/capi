// dev provides the development build, which may have unimplemented features,
// or dragons.. Who knows...
package main

import (
	// Standard Library Imports
	"fmt"
	"os"

	// External Imports
	log "github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/aciddude/capi"
	"github.com/aciddude/capi/datastore"
)

var (
	configPath = "../../config/config.yaml"
)

func main() {
	logger := log.WithFields(log.Fields{
		"app": "capi-dev",
	})

	config, err := capi.NewConfig(configPath)
	if err != nil {
		logger.WithError(err).Error("error processing config")
		os.Exit(1)
	}

	store, err := datastore.NewDatastore(config)
	if err != nil {
		logger.WithError(err).Error("error starting datastore")
		os.Exit(1)
	}

	fmt.Printf("%#+v\n", config)

	// TODO: Bind store into service.
	// TODO: configure API endpoints based on config.

	err = store.Close()
	if err != nil {
		logger.WithError(err).Error("error closing datastore connections")
		os.Exit(1)
	}

	// Close happily
	os.Exit(0)
}
