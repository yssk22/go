// Package uuid to generate UUID string
// We support only RFC4122 compatible one.
package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

// UUID is an alias for uuid string
type UUID [16]byte

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// New generates uuid v4 following RFC4122 Section 4.4
func New() UUID {
	uuid := make([]byte, 16, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(err) // should never fail in rand
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80 // (4.1.1) 10xxxxxx, The variant specified in this document.
	uuid[6] = uuid[6]&^0xf0 | 0x40 // (4.1.3) 0100xxxx, set version 4
	var _uuid [16]byte
	copy(_uuid[:], uuid)
	return UUID(_uuid)
}
