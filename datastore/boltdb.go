package datastore

import (
	// Standard Library Imports
	"os"
	"time"

	// External Imports
	bolt "go.etcd.io/bbolt"

	// Internal Imports
	"github.com/aciddude/capi"
)

const (
	boltdbPathBlocks       = "capi_blocks.db"
	boltdbPathTransactions = "capi_transactions.db"
)

// NewBoltDB returns a new BoltDB backed datastore.
func NewBoltDB(config *capi.Config) (*Datastore, error) {
	boltOptions := &bolt.Options{
		Timeout: time.Second * time.Duration(config.Datastore.BoltDB.Timeout),
	}

	// Get a listing of each coin in order to create boltDB .
	var coins []string
	for _, configuredCoin := range config.Coins {
		if configuredCoin.Code != "" {
			var found bool
			for _, coin := range coins {
				if coin == configuredCoin.Code {
					found = true
					break
				}
			}

			if !found {
				coins = append(coins, configuredCoin.Code)
			}
		}
	}

	// Ensure DB path exists, if not, create the path.
	ok, _ := fileExists(config.Datastore.BoltDB.DbPath)
	if !ok {
		err := os.MkdirAll(config.Datastore.BoltDB.DbPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	blocks := &Blocks{
		DBPath: config.Datastore.BoltDB.DbPath,

		DB:      nil,
		Coins:   coins,
		Options: boltOptions,
	}

	transactions := &Transactions{
		DBPath: config.Datastore.BoltDB.DbPath,

		DB:      nil,
		Coins:   coins,
		Options: boltOptions,
	}

	store := &Datastore{
		Blocks:       blocks,
		Transactions: transactions,
	}

	return store, nil
}
