// Package mongo implements `iam-pump/internal/pumpserver/store.Store` interface.
package mongo

import (
	"context"
	"sync"

	iamgo "github.com/che-kwas/iam-kit/mongo"
	"go.mongodb.org/mongo-driver/mongo"

	"iam-pump/internal/pumpserver/store"
)

type datastore struct {
	mgo *mongo.Client
}

func (ds *datastore) Audit() store.AuditStore {
	return newAudit(ds)
}

func (ds *datastore) Close(ctx context.Context) error {
	return ds.mgo.Disconnect(ctx)
}

var (
	mgoStore store.Store
	once     sync.Once
)

// MongoStore returns a mongo store instance.
func MongoStore() (store.Store, error) {
	if mgoStore != nil {
		return mgoStore, nil
	}

	var err error
	var mgo *mongo.Client
	once.Do(func() {
		mgo, err = iamgo.NewMongoIns()
		mgoStore = &datastore{mgo}
	})

	return mgoStore, err
}
