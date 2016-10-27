package example

// Example is an example of datastore model.
//go:generate dsmodel -type=Example
type Example struct {
	ID string `json:"id" ent:"id"`
}

type AliasNotUsed int

// no target
type A struct {
}
