package xos

import (
	"os"
)

// TryStats is like os.Stat but try multiple paths in order and returns the first found.
func TryStats(paths ...string) (os.FileInfo, error) {
	var err error
	var info os.FileInfo
	for _, p := range paths {
		info, err = os.Stat(p)
		if err == nil {
			return info, nil
		}
	}
	return nil, err
}
