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

// // GetterList is a list of Getter to fallback if
// type GetterList struct {
//     list []Getter
// }

// // Get implements `Getter.Get(string)` to try getting config value from the head of list.
// // If a Getter item returns the value, it would be returned.
// // If a Getter item returns an error other than KeyError, it fails and return that error immediately.
// func (c *GetterList) Get(key string) (interface{}, error) {
// 	for _, getter := range c.list {
// 		v, e := getter.Get(key)
// 		if e != nil {
// 			if _, ok := e.(KeyError); !ok {
// 				return nil, e
// 			}
// 		}
// 		if v != nil {
// 			return v, nil
// 		}
// 	}
// 	return nil, KeyError(key)
// }

// // GetStringOr is like Get and returns a string value or default value if not found.
// func (c *GetterList) GetStringOr(key string, or string) (string) {
//     v, e := c.Get(key)
//     if e != nil {
//         return or
//     }
//     return fmt.Sprintf("%s", v)
// }

// // GetStringOr is like Get and returns a string value or default value if not found.
// func (c *GetterList) GetIntOr(key string, or string) (string) {
//     v, e := c.Get(key)
//     if e != nil {
//         return or
//     }
//     return fmt.Sprintf("%s", v)
// }

// // Add adds a Getter to the list
// func (c *GetterList) Add(g Getter) {
//     c.list = append(c.list, g)
// }

// var defaultGetterList = &GetterList{}

// func Get(key string) (interface{}, error) {
//     return defaultGetterList.Get(key)
// }

// func AddGetter(g Getter) {
//     defaultGetterList.Add(g)
// }
