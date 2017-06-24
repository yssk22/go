package xbytes

import (
	"encoding/json"
)

// ByteString is an alias for []byte to encode/decode JSON as unicode string
type ByteString []byte

// MarshalJSON is to implement JSON marshalizer
func (s *ByteString) MarshalJSON() ([]byte, error) {
	bytes, err := json.Marshal(string(*s))
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalJSON is to implement JSON unmarshalizer
func (s *ByteString) UnmarshalJSON(data []byte) error {
	var x string
	err := json.Unmarshal(data, &x)
	*s = ByteString(x)
	return err
}
