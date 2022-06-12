package store

import (
	"context"
)

// AuditStore defines the audit storage interface.
type AuditStore interface {
	InsertMany(ctx context.Context, records []interface{}) error
}
