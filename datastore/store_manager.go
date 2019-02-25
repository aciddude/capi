package datastore

import "context"

// Manager brings together the interfaces required to manage a datastore.
type Manager interface {
	Creater
	Configurer
	Stater
}

// Creater provides an interface for any required setup that may be needed on
// first initialization for a given datastore.
type Creater interface {
	// CreateSchema provides a hook for the underlying storage to perform
	// operations in order to create the required store for a given resource.
	// coins is provided as a parameter in case individual tables are required
	// to be created per coin.
	CreateSchema(ctx context.Context, coins []string) error
}

// Configurer provides an interface for required configuration that may be
// required on startup.
// For example, ensuring indices are configured e.t.c.
type Configurer interface {
	// Configure provides a hook to configure the database.
	// coins is provided as a parameter in case individual tables require being
	// configured per coin.
	Configure(ctx context.Context, coins []string) error
}

// Stater provides simple state management for the datastore in order to
// provide state information based on coin.
type Stater interface {
	// IsCreated returns a bool to trigger crea
	IsCreated(ctx context.Context, coins []string) bool
	// Close enables datastore connections to be closed.
	Close() error
}
