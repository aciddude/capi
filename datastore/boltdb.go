package datastore

import (
	// Standard Library Imports
	"time"

	// External Imports
	bolt "go.etcd.io/bbolt"

	// Internal Imports
	"github.com/aciddude/capi"
)

// NewBoltDB returns a new BoltDB backed datastore.
func NewBoltDB(config *capi.Config) (*Datastore, error) {
	boltOptions := &bolt.Options{
		Timeout: time.Second * time.Duration(config.Datastore.BoltDB.Timeout),
	}

	blocksDb, err := bolt.Open("blocks.db", 0666, boltOptions)
	if err != nil {
		return nil, err
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

	blocks := &Blocks{
		DB:      blocksDb,
		Coins:   coins,
		Options: boltOptions,
	}

	transactionsDb, err := bolt.Open("transactions.db", 0666, boltOptions)
	if err != nil {
		return nil, err
	}

	transactions := &Transactions{
		DB:      transactionsDb,
		Coins:   coins,
		Options: boltOptions,
	}

	store := &Datastore{
		Blocks:       blocks,
		Transactions: transactions,
	}

	return store, nil
}
