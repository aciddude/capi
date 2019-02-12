package capi

import (
	// Standard Library Imports
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	// External Imports
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type datastore string

const (
	// BoltDB codifies the value to specify to use boltdb as a datastore.
	BoltDB datastore = "boltdb"
)

// Config provides the domain structure to enable configuring CAPI.
type Config struct {
	// Port configures the API port CAPI will serve from.
	Port string

	// Coin Configuration.
	Coins []Coin

	// Datastore configuration.
	Datastore Datastore
}

// Coin provides the configuration required to connect to a coin daemon API.
type Coin struct {
	// Name is the human readable name of the coin. For example, "Feathercoin".
	Name string
	// Code is the 3 letter coin code. For example, "FTC".
	Code string
	// Host is the Coin's API daemon hostname on which to connect to the API.
	Host string
	// Port is the Coin's API daemon port on which to connect to the API.
	Port string
	// Username is the username to use in order to authenticate to the Coin's
	// API daemon.
	Username string
	// Password is the password to use in order to authenticate to the Coin's
	// API daemon.
	Password string
	// Timeout is how long to wait before timing out API requests.
	Timeout int
	// SSL is whether to connect over SSL.
	// If not specified, the default is false.
	SSL bool
	// EnableCoinCodexAPI is whether to enable the coin's codex API.
	// If not specified, the default is false.
	EnableCoinCodexAPI bool `yaml:"enableCoinCodexAPI"`
}

// Datastore provides customization for backend datastores.
type Datastore struct {
	// Backend is the specific datastore driver to use.
	Backend datastore

	// BoltDB specific datastore configuration.
	BoltDB ConfigBoltDB `yaml:"boltdb"`
}

// ConfigBoltDB provides specific configuration customisation for BoltDB.
type ConfigBoltDB struct {
	// Timeout specifies, in seconds, how long to wait before timing out when
	// trying to gain a file lock on the database's BoltDB data file.
	Timeout int
}

// NewConfig returns a processed config object.
func NewConfig(configFilepath string) (*Config, error) {
	logger := log.WithFields(log.Fields{
		"package":  "capi",
		"function": "NewConfig",
	})

	if configFilepath == "" {
		return nil, errors.New("filepath is required to open config")
	}

	rawConfigFile, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		logger.WithError(err).Debug("filepath not provided")
		return nil, err
	}

	config := &Config{}
	switch {
	case strings.Contains(configFilepath, ".yaml"):
		err := yaml.Unmarshal(rawConfigFile, config)
		if err != nil {
			logger.WithError(err).Debug("error unmarshaling yaml config")
			return nil, err
		}

	case strings.Contains(configFilepath, ".json"):
		err := json.Unmarshal(rawConfigFile, config)
		if err != nil {
			logger.WithError(err).Debug("error unmarshaling json config")
			return nil, err
		}

	default:
		err := errors.New(fmt.Sprintf("'%s' contains an unprocessible config filetype", configFilepath))
		logger.WithError(err).Debug("no filetype found")
		return nil, err
	}

	return config, nil
}
