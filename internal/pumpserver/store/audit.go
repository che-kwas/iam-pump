package store

import (
	"context"
)

// AuditStore defines the audit storage interface.
type AuditStore interface {
	Create(ctx context.Context) error
}
