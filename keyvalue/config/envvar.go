package config

import (
	"fmt"
	"os"
	"unicode"

	"github.com/yssk22/go/keyvalue"
)

// EnvVar implements keyvalue.Getter for environment variables.
// It converts the keys into the actual environment variable names by
// following convensions.
//
//   - Every single char is converted to the upper case.
//      - key `abc` is converted as `ABC` environment variable name.
//   - non-alphabetical chars nor non-digits are converted to '_'
//      - key `abc.c` is `ABC_C`
//      - key `abc-bar` is `ABC_BAR`
//
var EnvVar = keyvalue.NewList(&envVar{})

type envVar struct {
}

func (e envVar) Get(key interface{}) (interface{}, error) {
	var skey string
	var ok bool
	if skey, ok = key.(string); !ok {
		return nil, keyvalue.KeyError(fmt.Sprintf("%s (not string)", key))
	}
	varname := getEnvVarName(skey)
	if v, ok := os.LookupEnv(varname); ok {
		return v, nil
	}
	return nil, keyvalue.KeyError(
		fmt.Sprintf(
			"%s (%s environment variable)",
			skey, varname,
		),
	)
}

func getEnvVarName(key string) string {
	u := []rune(key)
	for i, code := range u {
		if 'A' <= code && code <= 'Z' {
			// Nothing
		} else if 'a' <= code && code <= 'z' {
			u[i] = unicode.ToUpper(code)
		} else if '0' < code && code < '9' {
			// Nothing
		} else {
			u[i] = '_'
		}
	}
	return string(u)
}
