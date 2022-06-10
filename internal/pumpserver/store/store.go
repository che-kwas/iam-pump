// Package store defines the storage interface for iam-pump.
package store

import "context"

//go:generate mockgen -self_package=iam-pump/internal/pumpserver/store -destination mock_store.go -package store iam-pump/internal/pumpserver/store Factory,AuditStore

var client Store

// Store defines the pumpserver storage interface.
type Store interface {
	Audit() AuditStore
	Close(ctx context.Context) error
}

// Client returns the store client.
func Client() Store {
	return client
}

// SetClient sets the store client.
func SetClient(store Store) {
	client = store
}
