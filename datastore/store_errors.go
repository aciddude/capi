package datastore

import (
	"errors"
)

var (
	// ErrNotFound is the error returned for when a resource could not be
	// found in the datastore.
	ErrNotFound = errors.New("not found")
)
