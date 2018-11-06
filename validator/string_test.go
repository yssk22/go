package validator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_String_Min(t *testing.T) {
	a := assert.New(t)
	a.NotNil(String().Min(1).Validate(""))
	a.Nil(String().Min(1).Validate("s"))
	a.Nil(String().Min(1).Validate("ss"))
}

func Test_String_Max(t *testing.T) {
	a := assert.New(t)
	a.Nil(String().Max(1).Validate(""))
	a.Nil(String().Max(1).Validate("s"))
	a.NotNil(String().Max(1).Validate("ss"))
}

func Test_String_Range(t *testing.T) {
	a := assert.New(t)
	a.NotNil(String().Range(1, 2).Validate(""))
	a.Nil(String().Range(1, 2).Validate("s"))
	a.Nil(String().Range(1, 2).Validate("ss"))
	a.NotNil(String().Range(1, 2).Validate("sss"))
}
