// Copyright (C) 2015 SPEEDLAND Project
package assert

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
	"x/value"
)

// Assert is a helper for testing assersion
// You can use assert functions by casting *testing.T to Assert
//
//    import (
//        "testing"
//        "x/testing/assert"
//    )
//
//    func TestSomething(t *testing.T){
//        assert := assert.New(t)
//        assert.OK(true)
//        assert.EqInt(1, 1)
//        assert.EqStr("Expects", "Got", "Somthing wrong!")
//    }
//
type Assert struct {
	*testing.T
}

// New return a new *Assert
func New(t *testing.T) *Assert {
	return &Assert{t}
}

// OK for true assertion.
func (a *Assert) OK(ok bool, msgContext ...interface{}) {
	if !ok {
		a.failure("true", ok, msgContext...)
	}
}

// Not for false assertion
func (a *Assert) Not(ok bool, msgContext ...interface{}) {
	if ok {
		a.failure("true", ok, msgContext...)
	}
}

// Nil for nil assertion
func (a *Assert) Nil(v interface{}, msgContext ...interface{}) {
	if !value.IsNil(v) {
		a.failure("<nil>", v, msgContext...)
	}
}

// NotNil for non-nil assertion
func (a *Assert) NotNil(v interface{}, msgContext ...interface{}) {
	if value.IsNil(v) {
		a.failure("<non-nil>", v, msgContext...)
	}
}

// Zero for Zero value assertion
func (a *Assert) Zero(v interface{}, msgContext ...interface{}) {
	if !value.IsZero(v) {
		a.failure("<zero>", v, msgContext...)
	}
}

// NotZero for non Zero assertion
func (a *Assert) NotZero(v interface{}, msgContext ...interface{}) {
	if value.IsZero(v) {
		a.failure("<not-zero>", v, msgContext...)
	}
}

// EqInt for equality assertion (int)
func (a *Assert) EqInt(expect, got int, msgContext ...interface{}) {
	if expect != got {
		a.failure(expect, got, msgContext...)
	}
}

// EqInt32 for equality assertion (int32)
func (a *Assert) EqInt32(expect, got int32, msgContext ...interface{}) {
	if expect != got {
		a.failure(expect, got, msgContext...)
	}
}

// EqInt64 for equality assertion (int64)
func (a *Assert) EqInt64(expect, got int64, msgContext ...interface{}) {
	if expect != got {
		a.failure(expect, got, msgContext...)
	}
}

// EqFloat32 for equality assertion (float32)
func (a *Assert) EqFloat32(expect, got float32, msgContext ...interface{}) {
	if expect != got {
		a.failure(expect, got, msgContext...)
	}
}

// EqFloat64 for equality assertion (float64)
func (a *Assert) EqFloat64(expect, got float64, msgContext ...interface{}) {
	if expect != got {
		a.failure(expect, got, msgContext...)
	}
}

// EqStr for equality assertion (string)
func (a *Assert) EqStr(expect, got string, msgContext ...interface{}) {
	if expect != got {
		a.failure(expect, got, msgContext...)
	}
}

// EqByteString for equality assertion ([]byte string)
func (a *Assert) EqByteString(expect string, got []byte, msgContext ...interface{}) {
	if expect != string(got) {
		a.failure(expect, got, msgContext...)
	}
}

// EqTime for equality assertion (time.Time)
func (a *Assert) EqTime(expect, got time.Time, msgContext ...interface{}) {
	if !expect.Equal(got) {
		a.failure(expect, got, msgContext...)
	}
}

// GtInt for 'greater than' assertion (int)
func (a *Assert) GtInt(min, got int, msgContext ...interface{}) {
	if min > got {
		a.failure(fmt.Sprintf("> %d", min), got, msgContext...)
	}
}

// LtInt for 'less than' assertion (int)
func (a *Assert) LtInt(max, got int, msgContext ...interface{}) {
	if max < got {
		a.failure(fmt.Sprintf("< %d", max), got, msgContext...)
	}
}

func (a *Assert) failure(expect interface{}, got interface{}, msgContext ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	testpc, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	fun := runtime.FuncForPC(pc)
	packagepath := strings.Split(runtime.FuncForPC(testpc).Name(), ".")[0]

	report := fmt.Sprintf(
		"%s failure\n"+
			"\t source: %s/%s:%d\n"+
			"\t expect: %#v\n"+
			"\t    got: %#v",
		fun.Name(), packagepath, file, line, expect, got,
	)

	if len(msgContext) > 0 {
		a.Fatalf("%s\n\tcontext: %s", report, fmt.Sprintf(fmt.Sprintf("%s", msgContext[0]), msgContext[1:]...))
	} else {
		a.Fatalf("%s", report)
	}
}
