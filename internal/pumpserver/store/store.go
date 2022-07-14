// Package store defines the Store interface.
package store

import "context"

//go:generate mockgen -self_package=iam-pump/internal/pumpserver/store -destination mock_store.go -package store iam-pump/internal/pumpserver/store Store

var client Store

// Store defines the behavior of a store.
type Store interface {
	InsertOne(ctx context.Context, record interface{}) error
	InsertMany(ctx context.Context, records []interface{}) error
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
