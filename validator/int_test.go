package validator

import (
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_Int_Min(t *testing.T) {
	a := assert.New(t)
	a.NotNil(Int().Min(1).Validate(0))
	a.Nil(Int().Min(1).Validate(1))
	a.Nil(Int().Min(1).Validate(2))
}

func Test_Int_Max(t *testing.T) {
	a := assert.New(t)
	a.Nil(Int().Max(1).Validate(0))
	a.Nil(Int().Max(1).Validate(1))
	a.NotNil(Int().Max(1).Validate(2))
}

func Test_Int_Range(t *testing.T) {
	a := assert.New(t)
	a.NotNil(Int().Range(1, 2).Validate(0))
	a.Nil(Int().Range(1, 2).Validate(1))
	a.Nil(Int().Range(1, 2).Validate(2))
	a.NotNil(Int().Range(1, 2).Validate(3))
}
