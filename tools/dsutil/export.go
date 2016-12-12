package dsutil

import (
	"encoding/json"
	"io"
	"log"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
)

type ExportOption struct {
}

var DefaultExportOption = &ExportOption{}

const batchSize = 200

func Export(ctx context.Context, kind string, w io.Writer, option *ExportOption) error {
	if option == nil {
		option = DefaultExportOption
	}
	q := datastore.NewQuery(kind).Limit(batchSize)
	iter := q.Run(ctx)
	loop := 0
	totalCount := 0
	for {
		for j := 0; j < batchSize; j++ {
			var data entity
			var buff []byte
			data = entity(make(map[string]interface{}))
			key, err := iter.Next(data)
			if err != nil {
				if err == datastore.Done {
					log.Printf("Exported %d records", totalCount)
					return nil
				}
				return err
			}
			buff, err = json.Marshal(Row{
				Key: key, Data: data,
			})
			if err != nil {
				return err
			}
			w.Write(buff)
			w.Write([]byte("\n"))
			totalCount++
		}
		log.Printf("%d finished", (loop+1)*batchSize)
		cursor, err := iter.Cursor()
		if err != nil {
			return err
		}
		loop++
		iter = datastore.NewQuery(kind).Start(cursor).Limit(batchSize).Run(ctx)
	}
}
