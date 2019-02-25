package datastore

import (
	// Standard Library Imports
	"context"
	"errors"
	"fmt"
	"strings"

	// External Imports
	log "github.com/sirupsen/logrus"

	// Internal Imports
	"github.com/aciddude/capi"
)

// Datastore provides the concrete implementation of a datastore.
type Datastore struct {
	coins    []string
	managers []Manager

	Blocks       BlockManager
	Transactions TransactionManager
}

// CreateSchemas checks that each resource is created, if not, creates it.
func (s *Datastore) CreateSchemas(ctx context.Context) error {
	for _, manager := range s.managers {
		if manager == nil {
			continue
		}

		if !manager.IsCreated(ctx, s.coins) {
			err := manager.CreateSchema(ctx, s.coins)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ConfigureSchemas ensures that on start up, each resource is configured.
func (s *Datastore) ConfigureSchemas(ctx context.Context) error {
	for _, manager := range s.managers {
		err := manager.Configure(ctx, s.coins)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Datastore) Close() error {
	for _, manager := range s.managers {
		err := manager.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// NewDatastore returns a configured datastore.
func NewDatastore(config *capi.Config) (store *Datastore, err error) {
	logger := log.WithFields(log.Fields{
		"package":  "datastore",
		"function": "NewDatastore",
	})

	// context enables, where required, a datastore connection to be pushed
	// through the required setup procedures.
	var ctx context.Context

	// Choose the datastore backing based on user config.
	switch config.Datastore.Backend {
	case capi.BoltDB:
		store, err = NewBoltDB(config)

	default:
		err := errors.New(fmt.Sprintf("datastore '%s' is not implemented", config.Datastore.Backend))
		logger.WithError(err).Debug("datastore type not matched")
		return nil, err
	}

	// Bind in coins for database schema access.
	var coins []string
	for _, coin := range config.Coins {
		if coin.Code != "" {
			coins = append(coins, strings.ToLower(coin.Code))
		}
	}
	store.coins = coins

	// Bind datastore resource handlers into store.managers to enable easy
	// calling of functionality between resources.
	store.managers = []Manager{
		store.Blocks,
		store.Transactions,
	}

	err = store.CreateSchemas(ctx)
	if err != nil {
		logger.WithError(err).Debug("error creating datastore schemas")
		return nil, err
	}

	err = store.ConfigureSchemas(ctx)
	if err != nil {
		logger.WithError(err).Debug("error configuring datastore schemas")
		return nil, err
	}

	return store, nil
}
