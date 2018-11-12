package types

import (
	"encoding/json"

	"github.com/yssk22/go/x/xerrors"
)

// Typed to convert the non-typed interface i to typed interface v
func Typed(src interface{}, dst interface{}) {
	bytes, err := json.Marshal(src)
	xerrors.MustNil(err)
	xerrors.MustNil(json.Unmarshal(bytes, dst))
}
