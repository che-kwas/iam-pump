package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type audit struct {
	mgo *mongo.Client
}

func newAudit(ds *datastore) *audit {
	return &audit{ds.mgo}
}

// Create creates a new audit.
func (u *audit) Create(ctx context.Context) error {

	return nil
}
