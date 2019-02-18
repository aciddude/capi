package datastore

import (
	// Standard Library Imports
	"context"

	// Internal Imports
	"github.com/aciddude/capi"
)

// BlockManager pulls together all interfaces required to implement a database
// storage driver for storing coin blocks.
type BlockManager interface {
	Creater
	Configurer
	Stater
	BlockStorer
}

// BlockManager enables datastore implementers to filter for, create, update
// and delete blockchain block resources for multiple coins.
type BlockStorer interface {
	// List enables filtering for blocks in the datastore.
	List(ctx context.Context, coin string, query ListBlocksRequest) ([]*capi.Block, error)
	// Create creating a block in the datastore.
	Create(ctx context.Context, coin string, newEntity *capi.Block) (*capi.Block, error)
	// CreateBulk enables bulk creation of blocks in the datastore.
	CreateBulk(ctx context.Context, coin string, query []*capi.Block) ([]*capi.Block, error)
	// Get enables retrieving a block given an ID.
	Get(ctx context.Context, coin string, id string) (block *capi.Block, err error)
	// Update enables updating a block resource in the datastore.
	Update(ctx context.Context, coin string, id string, updatedEntity *capi.Block) (*capi.Block, error)
	// Delete removes a block from the datastore.
	Delete(ctx context.Context, coin string, id string) (err error)
}

// ListBlocksRequest enables filtering blocks based on the following
// parameters.
type ListBlocksRequest struct {
	// WalletID enables filtering blocks based on the a Wallet.
	WalletID string `json:"walletId"`
	// TransactionID enables filtering blocks based on a Transaction.
	TransactionID string `json:"TransactionId"`
}
