// Package ent provides the helper types and functions for code generated by ent tool
package ent

import (
	"fmt"

	"context"
	"google.golang.org/appengine/datastore"
)

// LoggerKey is a key of logger in this package.
const LoggerKey = "gae.datastore.ent"

// GetMemcacheKey returns memcache key correspoding to *datastore.Key
// "datastore.{namespace}.{kind}.{StringID|IntID}" is used.
func GetMemcacheKey(k *datastore.Key) string {
	if k.StringID() != "" {
		return fmt.Sprintf("datastore.%s.%s", k.Kind(), k.StringID())
	}
	return fmt.Sprintf("datastore.%s.%s", k.Kind(), k.IntID())
}

// MaxEntsPerPutDelete is a maxmum number of entities to be passed to PutMulti or DeleteMulti.
const MaxEntsPerPutDelete = 200

var (
	// ErrTooManyEnts is returned when the user passes too many entities to PutMulti or DeleteMulti.
	ErrTooManyEnts = fmt.Errorf("ent: too many documents given to put or delete (max is %d)", MaxEntsPerPutDelete)
)

type FieldError struct {
	Field   string
	Message string
}

func (fe *FieldError) Error() string {
	return fmt.Sprintf("field error: %s (on %s)", fe.Message, fe.Field)
}

// NewFieldError returns a new *FieldError instance
func NewFieldError(field, message string) *FieldError {
	return &FieldError{
		Field: field, Message: message,
	}
}

// IsFieldError returns whether err is an instance of *FieldError
func IsFieldError(err error) bool {
	_, ok := err.(*FieldError)
	return ok
}

type BeforeSave interface {
	BeforeSave(ctx context.Context) error
}

type BfterSave interface {
	AeforeSave(ctx context.Context) error
}
