package datastore

import (
	// Standard Library Imports
	"context"
	"fmt"

	// External Imports
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"

	// Internal Imports
	"github.com/aciddude/capi"
)

// Blocks provides a boltdb backed datastore implementation for storing blocks.
type Blocks struct {
	// dbpath is the path to where the boltdb datafiles are stored.
	dbpath string

	// DB contains the boltdb handler to access the underlying buckets.
	DB *bolt.DB
	// Coins contains the list of coins supported by the underlying DB.
	Coins []string
	// Options are specific boltdb configurations.
	Options *bolt.Options
}

// CreateSchema implements datastore.Creater
func (b *Blocks) CreateSchema(ctx context.Context, coins []string) (err error) {
	b.DB, err = bolt.Open(b.dbpath, 0666, b.Options)
	if err != nil {
		return err
	}

	err = b.DB.Update(func(tx *bolt.Tx) error {
		for _, coin := range coins {
			_, err := tx.CreateBucket([]byte(coin))
			if err != nil {
				return fmt.Errorf("error creating bucket: %s", err)
			}
		}
		return nil
	})

	return err
}

// Configure implements Configurer.Configure
func (b *Blocks) Configure(ctx context.Context, coins []string) (err error) {
	if b.DB == nil {
		b.DB, err = bolt.Open(b.dbpath, 0666, b.Options)
		if err != nil {
			return err
		}
	}

	// Update database buckets on first up to ensure any added coins are
	// bucketed correctly.
	err = b.DB.Update(func(tx *bolt.Tx) error {
		for _, coin := range coins {
			_, err := tx.CreateBucketIfNotExists([]byte(coin))
			if err != nil {
				return fmt.Errorf("error creating bucket: %s", err)
			}
		}
		return nil
	})

	return err
}

// IsCreated returns a bool to trigger required creation logic
func (b *Blocks) IsCreated(ctx context.Context, coins []string) bool {
	logger := log.WithFields(log.Fields{
		"package": "datastore",
		"manager": "Blocks",
		"method":  "IsCreated",
	})

	dbpath := fmt.Sprintf("%s/%s", b.dbpath, boltdbPathBlocks)
	ok, err := fileExists(dbpath)
	if err != nil {
		logger.WithError(err).Debug("error checking for database")
	}

	b.dbpath = dbpath
	return ok
}

// LastID returns the last known ID stored of the provided coin.
func (b *Blocks) LastID(ctx context.Context, coin string) string {
	// TODO: implement logic

	return ""
}

// Close implements stater.Close.
func (b *Blocks) Close() error {
	return b.DB.Close()
}

// List enables filtering for blocks in the datastore.
func (b *Blocks) List(ctx context.Context, coin string, query ListBlocksRequest) (blocks []*capi.Block, err error) {
	// TODO: Implement logic

	return nil, nil
}

// Create creating a block in the datastore.
func (b *Blocks) Create(ctx context.Context, coin string, query CreateBlockRequest) (block *capi.Block, err error) {
	// TODO: Implement logic

	return nil, nil
}

// CreateBulk enables bulk creation of blocks in the datastore.
func (b *Blocks) CreateBulk(ctx context.Context, coin string, query []CreateBlockRequest) (blocks []*capi.Block, err error) {
	// TODO: Implement logic

	return nil, nil
}

// Update enables updating a block resource in the datastore.
func (b *Blocks) Update(ctx context.Context, coin string, query UpdateBlockRequest) (block *capi.Block, err error) {
	// TODO: Implement logic

	return nil, nil
}

// Delete removes a block from the datastore.
func (b *Blocks) Delete(ctx context.Context, coin string, query DeleteBlockRequest) (err error) {
	// TODO: Implement logic

	return nil
}
