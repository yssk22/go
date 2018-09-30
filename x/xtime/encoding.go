package xtime

import (
	"fmt"
	"strconv"
	"time"

	"github.com/yssk22/go/x/xerrors"
)

// Timestamp is a time.Time alias that uses timestamp serialization for json.
type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(*t).Unix())), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return xerrors.Wrap(err, "could not unmarshal value %q as xtime.Timestamp", string(data))
	}
	*t = Timestamp(time.Unix(val, 0))
	return nil
}
