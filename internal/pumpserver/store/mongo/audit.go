package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	db  = "iam_authz_audit"
	col = "audit-logs"
)

type audit struct {
	col *mongo.Collection
}

func newAudit(ds *datastore) *audit {
	return &audit{ds.mgo.Database(db).Collection(col)}
}

// InsertMany inserts multiple records into the collection.
func (u *audit) InsertMany(ctx context.Context, records []interface{}) error {
	_, err := u.col.InsertMany(ctx, records)

	return err
}
