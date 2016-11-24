package dsutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine/datastore"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/appengine"
	"google.golang.org/appengine/remote_api"
)

type ExportedRow struct {
	Key  *datastore.Key
	Data entityLoader
}

const batchSize = 300

func Export(ctx context.Context, kind string, w io.Writer) error {
	q := datastore.NewQuery(kind).Limit(batchSize)
	iter := q.Run(ctx)
	loop := 0
	for {
		for j := 0; j < batchSize; j++ {
			data := entityLoader(make(map[string]interface{}))
			key, err := iter.Next(data)
			if err != nil {
				if err == datastore.Done {
					return nil
				}
				return err
			}
			buff, err := json.Marshal(&ExportedRow{
				Key:  key,
				Data: data,
			})
			if err != nil {
				return err
			}
			w.Write(buff)
			w.Write([]byte("\n"))
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

type entityLoader map[string]interface{}

func (l entityLoader) Load(props []datastore.Property) error {
	for _, p := range props {
		obj := map[string]interface{}{
			"Value":   p.Value,
			"NoIndex": p.NoIndex,
		}
		if p.Multiple {
			if _, ok := l[p.Name]; !ok {
				l[p.Name] = make([]interface{}, 0)
			}
			l[p.Name] = append(l[p.Name].([]interface{}), obj)
		} else {
			l[p.Name] = obj
		}
	}
	return nil
}

func (l entityLoader) Save() ([]datastore.Property, error) {
	// unused
	return nil, nil
}

var remoteCtxScopes = []string{
	"https://www.googleapis.com/auth/appengine.apis",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/cloud-platform",
}

var ErrNoRemoteAPIKeyIsConfigured = fmt.Errorf("no GOOGLE_APPLICATION_CREDENTIALS environment key is configured")

func GetRemoteContext(host string, namespace string) (context.Context, error) {
	var hc *http.Client
	var err error
	const defaultCredentialEnvKey = "GOOGLE_APPLICATION_CREDENTIALS"
	if os.Getenv(defaultCredentialEnvKey) != "" {
		hc, err = google.DefaultClient(oauth2.NoContext, remoteCtxScopes...)
	} else {
		return nil, ErrNoRemoteAPIKeyIsConfigured
	}
	if err != nil {
		return nil, err
	}
	ctx, err := remote_api.NewRemoteContext(host, hc)
	if err != nil {
		return nil, err
	}
	if namespace != "" {
		return appengine.Namespace(ctx, namespace)
	}
	return ctx, nil
}
