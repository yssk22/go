package example

import "time"

// Example is an example of datastore model.
//go:generate dsmodel -type=Example
type Example struct {
	ID                  string    `json:"id" ent:"id"`
	Digit               int       `json:"digit" ent:"form,resetifmissing" default:"10"`
	Desc                string    `json:"desc" ent:"form" default:"This is defualt value"`
	ContentBytes        []byte    `json:"content_bytes" ent:"form"`
	SliceType           []string  `json:"slice_type" ent:"form"`
	SliceFloatType      []float64 `json:"slice_float_type" ent:"form"`
	BoolType            bool      `json:"bool_type" ent:"form"`
	FloatType           float64   `json:"float_type" ent:"form"`
	CreatedAt           time.Time `json:"created_at" default:"$now"`
	UpdatedAt           time.Time `json:"updated_at" ent:"timestamp"`
	BeforeSaveProcessed bool      `json:"before_save_processed"`
	AfterSaveProcessed  bool      `json:"after_save_processed"`
}

type AliasNotUsed int

// no target
type A struct {
}
