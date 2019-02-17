package datastore

import (
	// Standard Library Imports
	"context"
	"fmt"

	// External Imports
	"github.com/asdine/storm"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"

	// Internal Imports
	"github.com/aciddude/capi"
)

// Blocks provides a boltdb backed datastore implementation for storing blocks.
type Blocks struct {
	// dbpath is the path to where the boltdb datafiles are stored.
	DBPath string

	// DB contains the boltdb handler to access the underlying buckets.
	DB *bolt.DB
	// storm is a wrapper around BoltDB that provides higher level ORM based
	// operations used internally to easily perform queries.
	storm *storm.DB
	// Coins contains the list of coins supported by the underlying DB.
	Coins []string
	// Options are specific boltdb configurations.
	Options *bolt.Options
}

// CreateSchema implements datastore.Creater
func (b *Blocks) CreateSchema(ctx context.Context, coins []string) (err error) {
	b.DB, err = bolt.Open(b.DBPath, 0666, b.Options)
	if err != nil {
		return err
	}

	b.storm, err = storm.Open(b.DBPath, storm.UseDB(b.DB))
	if err != nil {
		return err
	}

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

// Configure implements Configurer.Configure
func (b *Blocks) Configure(ctx context.Context, coins []string) (err error) {
	if b.DB == nil {
		b.DB, err = bolt.Open(b.DBPath, 0666, b.Options)
		if err != nil {
			return err
		}
	}

	if b.storm == nil {
		b.storm, err = storm.Open(b.DBPath, storm.UseDB(b.DB))
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

	dbpath := fmt.Sprintf("%s/%s", b.DBPath, boltdbPathBlocks)
	ok, err := fileExists(dbpath)
	if err != nil {
		logger.WithError(err).Debug("error checking for database")
	}

	b.DBPath = dbpath
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
func (b *Blocks) Create(ctx context.Context, coin string, newEntity *capi.Block) (block *capi.Block, err error) {
	bucket := b.storm.From(coin)

	if newEntity.ID == "" {
		newEntity.ID = uuid.New().String()
	}

	err = bucket.Save(newEntity)
	if err != nil {
		return block, err
	}

	return newEntity, nil
}

// CreateBulk enables bulk creation of blocks in the datastore.
func (b *Blocks) CreateBulk(ctx context.Context, coin string, query []CreateBlockRequest) ([]*capi.Block, error) {
	// TODO: Implement logic

	return nil, nil
}

// Get enables retrieving a block given an ID.
func (b *Blocks) Get(ctx context.Context, coin string, id string) (block *capi.Block, err error) {
	bucket := b.storm.From(coin)

	block = &capi.Block{}
	err = bucket.One("ID", id, block)
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return nil, ErrNotFound

		default:
			// TODO: debug logging for non-captured errors
			return nil, err
		}
	}

	return
}

// Update enables updating a block resource in the datastore.
func (b *Blocks) Update(ctx context.Context, coin string, id string, updatedEntity *capi.Block) (block *capi.Block, err error) {
	bucket := b.storm.From(coin)

	err = bucket.Update(updatedEntity)
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return nil, ErrNotFound

		default:
			// TODO: debug logging for non-captured errors
			return nil, err
		}
	}

	return updatedEntity, err
}

// Delete removes a block from the datastore.
func (b *Blocks) Delete(ctx context.Context, coin string, id string) (err error) {
	bucket := b.storm.From(coin)

	block := &capi.Block{ID: id}
	err = bucket.DeleteStruct(block)
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			return ErrNotFound

		default:
			// TODO: debug logging for non-captured errors
			return err
		}
	}

	return
}
