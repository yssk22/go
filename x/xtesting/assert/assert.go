// Package assert provides the assersion feature for *testing.T
// You can use assert functions by `assert.New(t)`
//
//    import (
//        "testing"
//        "github.com/yssk22/go/x/xtesting/assert"
//    )
//
//    func TestSomething(t *testing.T){
//        assert := assert.New(t)
//        assert.OK(true)
//        assert.EqInt(1, 1)
//        assert.EqStr("Expects", "Got", "Somthing wrong!")
//    }
//
package assert

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yssk22/go/x/xreflect"
)

// Assert is a helper struct for testing assersion
type Assert struct {
	*testing.T
}

// New return a new *Assert
func New(t *testing.T) *Assert {
	return &Assert{t}
}

// SkipIfErr skip the test if an error occurrs.
func (a Assert) SkipIfErr(err error) {
	a.Helper()
	if err != nil {
		a.Skipf(err.Error())
	}
}

// OK for true assertion.
func (a *Assert) OK(ok bool, msgContext ...interface{}) {
	a.Helper()
	if !ok {
		a.Failure("true", ok, msgContext...)
	}
}

// Not for false assertion
func (a *Assert) Not(ok bool, msgContext ...interface{}) {
	a.Helper()
	if ok {
		a.Failure("true", ok, msgContext...)
	}
}

// Nil for nil assertion
func (a *Assert) Nil(v interface{}, msgContext ...interface{}) {
	a.Helper()
	if !xreflect.IsNil(v) {
		a.Failure("<nil>", v, msgContext...)
	}
}

// NotNil for non-nil assertion
func (a *Assert) NotNil(v interface{}, msgContext ...interface{}) {
	a.Helper()
	if xreflect.IsNil(v) {
		a.Failure("<non-nil>", v, msgContext...)
	}
}

// Zero for Zero value assertion
func (a *Assert) Zero(v interface{}, msgContext ...interface{}) {
	a.Helper()
	if !xreflect.IsZero(v) {
		a.Failure("<zero>", v, msgContext...)
	}
}

// NotZero for non Zero assertion
func (a *Assert) NotZero(v interface{}, msgContext ...interface{}) {
	a.Helper()
	if xreflect.IsZero(v) {
		a.Failure("<not-zero>", v, msgContext...)
	}
}

// EqInt for equality assertion (int)
func (a *Assert) EqInt(expect, got int, msgContext ...interface{}) {
	a.Helper()
	if expect != got {
		a.Failure(expect, got, msgContext...)
	}
}

// EqInt32 for equality assertion (int32)
func (a *Assert) EqInt32(expect, got int32, msgContext ...interface{}) {
	a.Helper()
	if expect != got {
		a.Failure(expect, got, msgContext...)
	}
}

// EqInt64 for equality assertion (int64)
func (a *Assert) EqInt64(expect, got int64, msgContext ...interface{}) {
	a.Helper()
	if expect != got {
		a.Failure(expect, got, msgContext...)
	}
}

// EqFloat32 for equality assertion (float32)
func (a *Assert) EqFloat32(expect, got float32, msgContext ...interface{}) {
	a.Helper()
	if expect != got {
		a.Failure(expect, got, msgContext...)
	}
}

// EqFloat64 for equality assertion (float64)
func (a *Assert) EqFloat64(expect, got float64, msgContext ...interface{}) {
	a.Helper()
	if expect != got {
		a.Failure(expect, got, msgContext...)
	}
}

// EqStr for equality assertion (string)
func (a *Assert) EqStr(expect, got string, msgContext ...interface{}) {
	a.Helper()
	if expect != got {
		if strings.IndexByte(got, '\n') >= 0 {
			dmp := diffmatchpatch.New()
			diffs := dmp.DiffMain(expect, got, false)
			var buff bytes.Buffer
			buff.WriteString("[DIFF]\n")
			for _, diff := range diffs {
				text := diff.Text
				text = strings.Replace(text, "\n", "\\n\n", -1)
				switch diff.Type {
				case diffmatchpatch.DiffInsert:
					_, _ = buff.WriteString("\x1b[41m")
					_, _ = buff.WriteString(text)
					_, _ = buff.WriteString("\x1b[0m")
				case diffmatchpatch.DiffDelete:
					_, _ = buff.WriteString("\x1b[44m")
					_, _ = buff.WriteString(text)
					_, _ = buff.WriteString("\x1b[0m")
				case diffmatchpatch.DiffEqual:
					_, _ = buff.WriteString(text)
				}
			}
			a.Failure(expect, got, buff.String())
			return
		}
		a.Failure(expect, got, msgContext...)
	}
}

// EqByteString for equality assertion ([]byte string)
func (a *Assert) EqByteString(expect string, got []byte, msgContext ...interface{}) {
	a.Helper()
	if expect != string(got) {
		a.Failure(expect, string(got), msgContext...)
	}
}

// EqTime for equality assertion (time.Time)
func (a *Assert) EqTime(expect, got time.Time, msgContext ...interface{}) {
	a.Helper()
	if !expect.Equal(got) {
		a.Failure(expect.Format(time.RFC3339), got.Format(time.RFC3339), msgContext...)
	}
}

// GtInt for 'greater than' assertion (int)
func (a *Assert) GtInt(min, got int, msgContext ...interface{}) {
	a.Helper()
	if min > got {
		a.Failure(fmt.Sprintf("> %d", min), got, msgContext...)
	}
}

// LtInt for 'less than' assertion (int)
func (a *Assert) LtInt(max, got int, msgContext ...interface{}) {
	a.Helper()
	if max < got {
		a.Failure(fmt.Sprintf("< %d", max), got, msgContext...)
	}
}

// Failure fails the test with a report
// This function expects to be used by another assertion package, not by test code.
func (a *Assert) Failure(expect interface{}, got interface{}, msgContext ...interface{}) {
	a.Helper()
	// pc, _, _, _ := runtime.Caller(1)
	// testpc, file, line, _ := runtime.Caller(2)
	// file = filepath.Base(file)
	// fun := runtime.FuncForPC(pc)
	// packagepath := strings.Split(runtime.FuncForPC(testpc).Name(), ".")[0]
	var report string
	if len(msgContext) > 0 {
		str := fmt.Sprintf("%s", msgContext[0])
		rest := msgContext[1:]
		report = fmt.Sprintf(
			"%s failure: %s\n"+
				"\t expect: %v\n"+
				"\t    got: %v",
			a.T.Name(),
			fmt.Sprintf(str, rest...),
			expect, got,
		)

	} else {
		report = fmt.Sprintf(
			"%s failure\n"+
				"\t expect: %v\n"+
				"\t    got: %v",
			a.T.Name(), expect, got,
		)
	}
	a.Fatalf("%s", report)
}
