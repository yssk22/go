// Package uuid to generate UUID string
// We support only RFC4122 compatible one.
package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
)

// UUID is an alias for uuid string
type UUID [16]byte

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// MarshalJSON implements json.Marshaler#MarshalJSON()
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", u.String())), nil
}

// UnmarshalJSON implements json.Unmarshaler#UnmarshalJSON()
func (u *UUID) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("Invalid string")
	}
	newval, ok := FromString(string(b[1 : len(b)-1]))
	if !ok {
		return fmt.Errorf("Invalid format")
	}
	*u = newval
	return nil
}

// New generates uuid v4 following RFC4122 Section 4.4
func New() UUID {
	_uuid := make([]byte, 16, 16)
	_, err := io.ReadFull(rand.Reader, _uuid)
	if err != nil {
		panic(err) // should never fail in rand
	}
	_uuid[8] = _uuid[8]&^0xc0 | 0x80 // (4.1.1) 10xxxxxx, The variant specified in this document.
	_uuid[6] = _uuid[6]&^0xf0 | 0x40 // (4.1.3) 0100xxxx, set version 4
	var uuid [16]byte
	copy(uuid[:], _uuid)
	return UUID(uuid)
}

var delimiter = []byte("-")[0]

// FromString returns *UUID from string
func FromString(s string) (UUID, bool) {
	if len(s) != 36 {
		return UUID{}, false
	}
	var uuid [16]byte
	// 012345678901234567890123456789012345
	// 99cf7a36-9ed7-4930-a992-f9966356329c
	var j = 0
	for i := 0; i < 16; i++ {
		var b int64
		var err error
		if s[j] == delimiter {
			j++
		}
		if b, err = strconv.ParseInt(s[j:j+2], 16, 16); err != nil {
			fmt.Printf("ParseError: %s - %v\n", s[j:j+2], err)
			return UUID{}, false
		}
		uuid[i] = byte(b)
		j += 2
	}
	return UUID(uuid), true
}
