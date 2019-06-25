package example

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/appengine"

	"github.com/yssk22/go/types"
)

// Entity is an example for datastore entity
// @datastore
type Entity struct {
	ID           string             `json:"id" ent:"key"`
	Digit        int                `json:"digit"`
	Desc         string             `json:"desc"`
	ContentBytes []byte             `json:"content_bytes" ent:"search"`
	SliceType    []string           `json:"slice_type"`
	BoolType     bool               `json:"bool_type" ent:"search"`
	FloatType    float64            `json:"float_type" ent:"search"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" ent:"timestamp"`
	CustomType   types.RGB          `json:"custom_type"`
	Location     appengine.GeoPoint `json:"location" ent:"search"`

	FieldWithNoIndex int `datastore:",noindex"`
	FieldWithHyphen  int `datastore:"-"`

	BeforeSaveDesc string
	AfterSaveDesc  string
}

func (e *Entity) BeforeSave(ctx context.Context) error {
	e.BeforeSaveDesc = fmt.Sprintf("(BeforeSave) %s", e.Desc)
	return nil
}

func (e *Entity) AfterSave(ctx context.Context) error {
	e.AfterSaveDesc = fmt.Sprintf("(AfterSave) %s", e.Desc)
	return nil
}
