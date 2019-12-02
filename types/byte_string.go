package types

import (
	"encoding/json"
)

// ByteString is an alias type of []byte to avoid base64 encoding in JSON.
type ByteString []byte

// MarshalJSON implements json.Marshaler#MarshalJSON()
func (bs ByteString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(bs))
}

// UnmarshalJSON implements json.Unmarshaler#UnmarshalJSON()
func (bs *ByteString) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*bs = []byte(s)
	return nil
}
