package datastore

import (
	// Standard Library Imports
	"context"
	"fmt"

	// External Imports
	bolt "go.etcd.io/bbolt"

	// Internal Imports
	"github.com/aciddude/capi"
)

// Blocks provides a boltdb backed datastore implementation for storing blocks.
type Blocks struct {
	DB *bolt.DB
	// the list of coins supported by the underlying DB.
	Coins []string
	// options are specific boltdb configurtions.
	Options *bolt.Options
}

// CreateSchema implements datastore.Creater
func (b *Blocks) CreateSchema(ctx context.Context, coins []string) error {
	for _, coin := range coins {
		err := b.DB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte(coin))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// Configure implements Configurer.Configure
func (b *Blocks) Configure(ctx context.Context, coins []string) error {
	// nothing to configure here...
	return nil
}

// IsCreated returns a bool to trigger crea
func (b *Blocks) IsCreated(ctx context.Context, coins []string) bool {
	// TODO: implement logic

	return false
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
