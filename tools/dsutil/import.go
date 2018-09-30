package dsutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"google.golang.org/appengine/datastore"

	"context"
)

type ImportOption struct {
	AppID string
	Skip  int
}

var DefaultImportOption = &ImportOption{}

func Import(ctx context.Context, kind string, r io.Reader, option *ImportOption) (int, error) {
	if option == nil {
		option = DefaultImportOption
	}
	var keys []*datastore.Key
	var ents []entity
	reader := bufio.NewReader(r)
	numImported := 0
	linesProcessed := 0
	for {
		buff, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return numImported, fmt.Errorf("io error: %v", err)
		}
		linesProcessed++
		if option.Skip > 0 && linesProcessed <= option.Skip {
			continue
		}
		var row Row
		if err := json.Unmarshal(buff, &row); err != nil {
			continue
		}
		key := datastore.NewKey(ctx, row.Key.Kind(), row.Key.StringID(), row.Key.IntID(), nil)
		keys = append(keys, key)
		ents = append(ents, row.Data)
		if len(keys) >= batchSize {
			if _, err := datastore.PutMulti(ctx, keys, ents); err != nil {
				return numImported, err
			}
			numImported += len(keys)
			log.Printf("%d rows imported", numImported)
			keys = make([]*datastore.Key, 0)
			ents = make([]entity, 0)
		}
	}
	if len(keys) > 0 {
		if _, err := datastore.PutMulti(ctx, keys, ents); err != nil {
			return numImported, err
		}
		numImported += len(keys)
		log.Printf("%d rows imported", numImported)
	}
	return numImported, nil
}
