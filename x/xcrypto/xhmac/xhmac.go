// Package xhmac provides some ulitities for hmac
package xhmac

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
)

// ErrWrongSignatureFormat is an error object when receiving the wrong signature.
var ErrWrongSignatureFormat = fmt.Errorf("xhmac: wrong signature format")

// ErrSignatureModified is an error object when the modified singed value is passed.
var ErrSignatureModified = fmt.Errorf("xhmac: signed value was modified")

var signatureDelimiter = []byte(".")

// Base64 is hmac utility to sign/unsign hmac with base64-encoding.
type Base64 struct {
	h   func() hash.Hash
	key []byte
}

// NewBase64 returns *Base64HMAC
func NewBase64(key []byte, h func() hash.Hash) *Base64 {
	if h == nil {
		h = sha256.New
	}
	return &Base64{
		h:   h,
		key: key,
	}
}

// Sign returns []byte where HMAC base64-encoded signature is appended.
func (b *Base64) Sign(value []byte) []byte {
	mac := b.genMAC(value)
	macb64Len := base64.StdEncoding.EncodedLen(len(mac))
	dst := make([]byte, macb64Len, macb64Len)
	base64.StdEncoding.Encode(dst, mac)
	result := append(value, signatureDelimiter...)
	result = append(result, dst...)
	return result
}

// SignString is a string version of Sign
func (b *Base64) SignString(value string) string {
	return string(b.Sign([]byte(value)))
}

// Unsign returns a value extracted from a signed message made by Sign.
func (b *Base64) Unsign(singedMessage []byte) ([]byte, error) {
	sep := bytes.LastIndex(singedMessage, signatureDelimiter)
	if sep < 0 {
		return nil, ErrWrongSignatureFormat
	}
	value := singedMessage[:sep]
	macb64 := singedMessage[sep+1:]
	macb64Len := base64.StdEncoding.DecodedLen(len(macb64))
	dst := make([]byte, macb64Len, macb64Len)

	n, err := base64.StdEncoding.Decode(dst, macb64)
	if err != nil {
		return nil, ErrWrongSignatureFormat
	}
	mac := dst[:n]
	if !hmac.Equal(mac, b.genMAC(value)) {
		return nil, ErrSignatureModified
	}
	return value, nil
}

// UnsignString is a string version of Unsign
func (b *Base64) UnsignString(value string) (string, error) {
	unsigned, err := b.Unsign([]byte(value))
	if err != nil {
		return "", err
	}
	return string(unsigned), nil
}

func (b *Base64) genMAC(value []byte) []byte {
	mac := hmac.New(b.h, b.key)
	mac.Write(value)
	return mac.Sum(nil)
}
