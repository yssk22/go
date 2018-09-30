// Package config provies configuraiton access functions
package config

import "github.com/yssk22/go/keyvalue"

var defaultList = keyvalue.NewList()

// Setup initialize the Getters for Get* funcitons
func Setup(g ...keyvalue.Getter) {
	defaultList = keyvalue.NewList(g...)
}

// Get returns a value from the default config list
func Get(key string) (interface{}, error) {
	return defaultList.Get(key)
}

// GetOr returns a value from the default config list or a default value if not found.
func GetOr(key string, or interface{}) interface{} {
	return defaultList.GetOr(key, or)
}

// GetStringOr is a string version of GetOr
func GetStringOr(key string, or string) string {
	return defaultList.GetStringOr(key, or)
}

// GetIntOr is a int version of GetOr
func GetIntOr(key string, or int) int {
	return defaultList.GetIntOr(key, or)
}
