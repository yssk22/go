package dsutil

import (
	"log"

	"github.com/yssk22/go/x/xtime"

	"google.golang.org/appengine/datastore"

	"time"

	"context"
)

type Iterator struct {
	BatchSize     int
	StopOnError   bool
	QueryInterval time.Duration
	q             *datastore.Query
}

// NewIterator returns a new *Iterator object
func NewIterator(q *datastore.Query) *Iterator {
	return &Iterator{
		BatchSize:     100,
		StopOnError:   true,
		QueryInterval: time.Duration(0),
		q:             q,
	}
}

// Run runs the iterator
func (i *Iterator) Run(ctx context.Context, gen func() interface{}, f func(v interface{}) error) error {
	q := i.q.Limit(i.BatchSize)
	var iter *datastore.Iterator
	t := xtime.Benchmark(func() {
		iter = q.Run(ctx)
	})
	log.Printf("Query takes %s", t)
	loop := 0
	totalCount := 0
	errorCount := 0
	for {
		for j := 0; j < i.BatchSize; j++ {
			record := gen()
			key, err := iter.Next(record)
			if err != nil {
				if err == datastore.Done {
					log.Printf("Processed %d records (error: %d)", totalCount, errorCount)
					return nil
				}
				if i.StopOnError {
					log.Printf("[error] %s: %v", key, err)
					return err
				}
				log.Printf("[warn] %v", err)
				errorCount++
			} else {
				if err := f(record); err != nil {
					if i.StopOnError {
						log.Printf("[error] %s: %v", key, err)
						return err
					}
					log.Printf("[warn] %v", err)
					errorCount++
				}
			}
			totalCount++
		}
		log.Printf("Processed %d records", (loop+1)*i.BatchSize)
		cursor, err := iter.Cursor()
		if err != nil {
			return err
		}
		loop++
		if int(i.QueryInterval) > 0 {
			time.Sleep(i.QueryInterval)
		}
		q = q.Start(cursor)
		t := xtime.Benchmark(func() {
			iter = q.Run(ctx)
		})
		log.Printf("Query takes %s", t)
	}
}
