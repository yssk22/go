package datastore

import (
	"fmt"

	"context"

	"google.golang.org/appengine/datastore"
)

// DeleteAll deletes the all `kind` entities stored in datastore
func DeleteAll(ctx context.Context, kind string) error {
	const batchSize = 300
	var keys []*datastore.Key
	var dummy []interface{}
	var err error
	for {
		if keys, err = datastore.NewQuery(kind).KeysOnly().Limit(batchSize).GetAll(ctx, dummy); err != nil {
			return fmt.Errorf("delete_all: error retrieving keys: %v", err)
		}
		if err := datastore.DeleteMulti(ctx, keys); err != nil {
			return fmt.Errorf("delete_all: error deleting keys: %v", err)
		}
		count, err := datastore.NewQuery(kind).KeysOnly().Count(ctx)
		if err != nil {
			return fmt.Errorf("delete_all: error checking remaining keys: %v", err)
		}
		if count == 0 {
			break
		}
	}
	return nil
}
