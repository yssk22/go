package example

import (
	"time"

	"google.golang.org/appengine"

	"github.com/yssk22/go/types"
)

// Entity is an example for datastore entity
// @datastore
type Entity struct {
	ID                  string             `json:"id" ent:"key"`
	Digit               int                `json:"digit"`
	Desc                string             `json:"desc"`
	ContentBytes        []byte             `json:"content_bytes" ent:"search"`
	SliceType           []string           `json:"slice_type"`
	BoolType            bool               `json:"bool_type" ent:"search"`
	FloatType           float64            `json:"float_type" ent:"search"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at" ent:"timestamp"`
	BeforeSaveProcessed bool               `json:"before_save_processed"`
	CustomType          types.RGB          `json:"custom_type"`
	Location            appengine.GeoPoint `json:"location" ent:"search"`

	FieldWithNoIndex int `datastore:",noindex"`
	FieldWithHyphen  int `datastore:"-"`
}
