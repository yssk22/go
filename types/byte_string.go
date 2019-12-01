package types

import "fmt"

// ByteString is an alias type of []byte to avoid base64 encoding in JSON.
type ByteString []byte

// MarshalJSON implements json.Marshaler#MarshalJSON()
func (bs ByteString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", string(bs))), nil
}

// UnmarshalJSON implements json.Unmarshaler#UnmarshalJSON()
func (bs *ByteString) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("Invalid string")
	}
	*bs = b[1 : len(b)-1]
	return nil
}
