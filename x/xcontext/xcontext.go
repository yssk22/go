package xcontext

import (
	"fmt"

	"github.com/speedland/go/x/xruntime"
)

// Key provides package unique context key
type Key struct {
	pkgName string
	key     string
}

// NewKey returns *Key
func NewKey(key string) *Key {
	return &Key{
		pkgName: xruntime.CaptureCaller().PackageName,
		key:     key,
	}
}

func (k *Key) String() string {
	return fmt.Sprintf("%s@%s", k.key, k.pkgName)
}
