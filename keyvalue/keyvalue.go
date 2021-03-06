// Package keyvalue provides the Key-Value based access utility.
package keyvalue

import (
	"fmt"
	"time"

	"github.com/yssk22/go/x/xtime"

	"github.com/yssk22/go/number"
)

// Getter is an interface to get a value by a key
type Getter interface {
	Get(interface{}) (interface{}, error)
}

// Setter is an interface to set a value bey a key
type Setter interface {
	Set(interface{}, interface{}) error
}

// GetterSetter is an interface to get/set a valeu by a key
type GetterSetter interface {
	Get(interface{}) (interface{}, error)
	Set(interface{}, interface{}) error
}

// KeyError is the error when the key is not found.
type KeyError string

func (e KeyError) Error() string {
	return fmt.Sprintf("key %q is not found", string(e))
}

// IsKeyError returns the error is KeyError or not
func IsKeyError(err error) bool {
	_, ok := err.(KeyError)
	return ok
}

// GetOr gets a value from Getter or return the defalut `or` value if not found.
func GetOr(g Getter, key interface{}, or interface{}) interface{} {
	if g == nil {
		return or
	}
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	return fmt.Sprintf("%s", v)
}

// GetStringOr is string version of GetOr.
func GetStringOr(g Getter, key interface{}, or string) string {
	if g == nil {
		return or
	}
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	switch v.(type) {
	case []string:
		// Take the first element follwoing to url.Values implementation
		l := v.([]string)
		if len(l) > 0 {
			return l[0]
		}
		return or
	default:
		return fmt.Sprintf("%s", v)
	}
}

// GetIntOr is int version of GetOr.
func GetIntOr(g Getter, key interface{}, or int) int {
	if g == nil {
		return or
	}
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

// GetFloatOr is float64 version of GetOr.
func GetFloatOr(g Getter, key interface{}, or float64) float64 {
	if g == nil {
		return or
	}
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	switch v.(type) {
	case int8:
		return float64(v.(int8))
	case int16:
		return float64(v.(int16))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case int:
		return float64(v.(int))
	case string:
		return number.ParseFloat64Or(v.(string), or)
	default:
		return or
	}
}

func GetDateOr(g Getter, key interface{}, or time.Time) time.Time {
	if g == nil {
		return or
	}
	v, e := g.Get(key)
	if e != nil {
		return or
	}
	switch v.(type) {
	case string:
		t, e := xtime.ParseDate(v.(string), or.Location(), or.Year())
		if e == nil {
			return t
		}
		return or
	default:
		return or
	}
}
