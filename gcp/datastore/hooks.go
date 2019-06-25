package datastore

import "context"

// BeforeSave is an interface to run a logic before save
type BeforeSave interface {
	BeforeSave(context.Context) error
}

// AfterSave is an interface to run a logic before save
type AfterSave interface {
	AfterSave(context.Context) error
}
