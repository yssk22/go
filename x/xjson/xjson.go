// Package xjson provides extended utility functions for encoding/json
package xjson

import (
	"encoding/json"
	"os"
)

// FromFile loads a JSON filepath into `v`
func FromFile(filepath string, v interface{}) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(v); err != nil {
		return err
	}
	return nil
}
