package example

import (
	"time"

	"github.com/speedland/go/rgb"
)

// Example is an example of datastore model.
//go:generate ent -type=Example
type Example struct {
	ID                  string    `json:"id" ent:"id"`
	Digit               int       `json:"digit" ent:"form,resetifmissing" default:"10"`
	Desc                string    `json:"desc" ent:"form" default:"This is default value"`
	ContentBytes        []byte    `json:"content_bytes" ent:"form"`
	SliceType           []string  `json:"slice_type" ent:"form"`
	BoolType            bool      `json:"bool_type" ent:"form"`
	FloatType           float64   `json:"float_type" ent:"form"`
	CreatedAt           time.Time `json:"created_at" default:"$now"`
	UpdatedAt           time.Time `json:"updated_at" ent:"timestamp"`
	DefaultTime         time.Time `json:"default_time" default:"2016-01-01T20:12:10Z"`
	BeforeSaveProcessed bool      `json:"before_save_processed"`
	AfterSaveProcessed  bool      `json:"after_save_processed"`
	CustomType          rgb.RGB   `json:"custom_type" ent:"form" parser:"github.com/speedland/go/rgb.MustParseRGB"`
}

type AliasNotUsed int

// no target
type A struct {
}
