package datastore

import (
	// Standard Library Imports
	"context"

	// Internal Imports
	"github.com/aciddude/capi"
)

// TransactionManager pulls together all interfaces required to implement a
// database storage driver for storing coin transactions.
type TransactionManager interface {
	Creater
	Configurer
	Stater
	TransactionStorer
}

// TransactionManager enables datastore implementers to filter for, create, update
// and delete blockchain transaction resources for multiple coins.
type TransactionStorer interface {
	// List enables filtering for transactions in the datastore.
	List(ctx context.Context, coin string, query ListTransactionsRequest) (transactions []*capi.Transaction, err error)
	// Create creating a transaction in the datastore.
	Create(ctx context.Context, coin string, query CreateTransactionRequest) (transaction *capi.Transaction, err error)
	// CreateBulk enables bulk creation of transactions in the datastore.
	CreateBulk(ctx context.Context, coin string, query []CreateTransactionRequest) (transactions []*capi.Transaction, err error)
	// Update enables updating a transaction resource in the datastore.
	Update(ctx context.Context, coin string, query UpdateTransactionRequest) (transaction *capi.Transaction, err error)
	// Delete removes a transaction from the datastore.
	Delete(ctx context.Context, coin string, query DeleteTransactionRequest) (err error)
}

// ListTransactionsRequest enables filtering transactions based on the following
// parameters.
type ListTransactionsRequest struct {
	// WalletID enables filtering transactions based on a Wallet.
	WalletID string `json:"walletId"`
}

// CreateTransactionRequest creates a Transaction.
type CreateTransactionRequest struct {
	Transaction *capi.Transaction
}

// UpdateTransactionRequest updates the transaction, with the newly provided resource.
type UpdateTransactionRequest struct {
	ID          string
	Transaction *capi.Transaction
}

// DeleteTransactionRequest deletes the specified transaction.
type DeleteTransactionRequest struct {
	ID string
}
