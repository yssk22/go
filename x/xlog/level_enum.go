// Code generated by "enum -type=Level"; DO NOT EDIT

package xlog

import (
	"encoding/json"
	"fmt"
)

var (
	_LevelValueToString = map[Level]string{
		LevelAll:   "all",
		LevelTrace: "trace",
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
		LevelFatal: "fatal",
		LevelNone:  "none",
	}
	_LevelStringToValue = map[string]Level{
		"all":   LevelAll,
		"trace": LevelTrace,
		"debug": LevelDebug,
		"info":  LevelInfo,
		"warn":  LevelWarn,
		"error": LevelError,
		"fatal": LevelFatal,
		"none":  LevelNone,
	}
)

func (i Level) String() string {
	if str, ok := _LevelValueToString[i]; ok {
		return str
	}
	return fmt.Sprintf("Level(%d)", i)
}

func ParseLevel(s string) (Level, error) {
	if val, ok := _LevelStringToValue[s]; ok {
		return val, nil
	}
	return Level(0), fmt.Errorf("Invalid value %q for Level", s)
}

func ParseLevelOr(s string, or Level) Level {
	val, err := ParseLevel(s)
	if err != nil {
		return or
	}
	return val
}

func (i Level) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _LevelValueToString[i]; !ok {
		s = fmt.Sprintf("Level(%d)", i)
	}
	return json.Marshal(s)
}

func (i *Level) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("Invalid string")
	}
	newval, err := ParseLevel(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*i = newval
	return nil
}