// Package mongo implements `iam-pump/internal/pumpserver/store.Store` interface.
package mongo

import (
	"context"

	iamgo "github.com/che-kwas/iam-kit/mongo"
	"go.mongodb.org/mongo-driver/mongo"

	"iam-pump/internal/pumpserver/store"
)

const (
	db  = "iam_authz_audit"
	col = "audit-logs"
)

type mgoStore struct {
	mgo *mongo.Client
}

var _ store.Store = &mgoStore{}

// InsertMany inserts multiple records into the collection.
func (m *mgoStore) InsertMany(ctx context.Context, records []interface{}) error {
	_, err := m.mgo.Database(db).Collection(col).InsertMany(ctx, records)

	return err
}

func (m *mgoStore) Close(ctx context.Context) error {
	return m.mgo.Disconnect(ctx)
}

// NewMongoStore returns a mongo store instance.
func NewMongoStore() (store.Store, error) {
	mgo, err := iamgo.NewMongoIns()
	if err != nil {
		return nil, err
	}

	return &mgoStore{mgo}, err
}
