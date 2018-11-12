package types

import (
	"encoding/json"
	"fmt"
)

var (
	_MyEnumValueToString = map[MyEnum]string{
		MyEnumFoo: "foo",
		MyEnumBar: "bar",
	}
	_MyEnumStringToValue = map[string]MyEnum{
		"foo": MyEnumFoo,
		"bar": MyEnumBar,
	}
)

func (e MyEnum) String() string {
	if str, ok := _MyEnumValueToString[e]; ok {
		return str
	}
	return fmt.Sprintf("MyEnum(%d)", e)
}

func (e MyEnum) IsVaild() bool {
	_, ok := _MyEnumValueToString[e]
	return ok
}

func ParseMyEnum(s string) (MyEnum, error) {
	if val, ok := _MyEnumStringToValue[s]; ok {
		return val, nil
	}
	return MyEnum(0), fmt.Errorf("invalid value %q for MyEnum", s)
}

func ParseMyEnumOr(s string, or MyEnum) MyEnum {
	val, err := ParseMyEnum(s)
	if err != nil {
		return or
	}
	return val
}

func MustParseMyEnum(s string) MyEnum {
	val, err := ParseMyEnum(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (e MyEnum) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _MyEnumValueToString[e]; !ok {
		s = fmt.Sprintf("MyEnum(%d)", e)
	}
	return json.Marshal(s)
}

func (e *MyEnum) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := ParseMyEnum(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*e = newval
	return nil
}

func (e *MyEnum) Parse(s string) error {
	if val, ok := _MyEnumStringToValue[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid value %q for MyEnum", s)
}
