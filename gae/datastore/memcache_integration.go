package datastore

import (
	"fmt"

	"google.golang.org/appengine/datastore"
)

// GetMemcacheKey returns memcache key correspoding to *datastore.Key
// "datastore.{namespace}.{kind}.{StringID|IntID}" is used.
func GetMemcacheKey(k *datastore.Key) string {
	if k.StringID() != "" {
		return fmt.Sprintf("datastore.%s.%s", k.Kind(), k.StringID())
	}
	return fmt.Sprintf("datastore.%s.%d", k.Kind(), k.IntID())
}
