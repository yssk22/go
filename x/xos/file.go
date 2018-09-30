package xos

import (
	"os"
)

// TryStats is like os.Stat but try multiple paths in order and returns the first found.
func TryStats(paths ...string) (os.FileInfo, error) {
	for _, p := range paths {
		if info, err := os.Stat(p); err == nil {
			return info, nil
		}
	}
	return nil, os.ErrNotExist
}

// TryExists is like TryStats but returns a path found.
func TryExists(paths ...string) (string, error) {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", os.ErrNotExist
}
