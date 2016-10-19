// Package keyvalue provides the Key-Value based access utility.
package keyvalue

import (
	"fmt"

	"github.com/speedland/go/number"
)

// Getter is an interface to get the config value
type Getter interface {
	Get(string) (interface{}, error)
}

// KeyError is the error when the key is not found.
type KeyError string

func (e KeyError) Error() string {
	return fmt.Sprintf("key %q is not found", string(e))
}

// Map is an alias for map[string]interface{} that implements Getter interface
type Map map[string]interface{}

// Get implements Getter#Get
func (m Map) Get(key string) (interface{}, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return nil, KeyError(key)
}

// GetgOr gets a value from Getter or return the defalut `or` value if not found.
func GetOr(g Getter, key string, or interface{}) interface{} {
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	return fmt.Sprintf("%s", v)
}

// GetStringOr is string version of GetOr.
func GetStringOr(g Getter, key string, or string) string {
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	return fmt.Sprintf("%s", v)
}

// GetIntOr is int version of GetOr.
func GetIntOr(g Getter, key string, or int) int {
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	switch v.(type) {
	case int8:
		return int(v.(int8))
	case int16:
		return int(v.(int16))
	case int32:
		return int(v.(int32))
	case int64:
		return int(v.(int64))
	case int:
		return v.(int)
	case string:
		return number.ParseIntOr(v.(string), or)
	default:
		return or
	}
}
