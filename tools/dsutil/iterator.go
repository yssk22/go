package dsutil

import (
	"log"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

func Iterator(ctx context.Context, q *datastore.Query, batchSize int, gen func() interface{}, f func(v interface{}) error) error {
	q = q.Limit(batchSize)
	iter := q.Run(ctx)
	loop := 0
	totalCount := 0
	for {
		for j := 0; j < batchSize; j++ {
			record := gen()
			key, err := iter.Next(record)
			if err != nil {
				if err == datastore.Done {
					log.Printf("Processed %d records", totalCount)
					return nil
				}
				return err
			}
			if err := f(record); err != nil {
				log.Println("[error] %s: %v", key, err)
				return err
			}
			totalCount++
		}
		log.Printf("Processed %d records", (loop+1)*batchSize)
		cursor, err := iter.Cursor()
		if err != nil {
			return err
		}
		loop++
		q = q.Start(cursor)
		iter = q.Run(ctx)
	}
}
