package example

import (
	"encoding/json"
	"fmt"
)

var (
	_MyEnumValueToString = map[MyEnum]string{
		MyEnumA: "a",
		MyEnumB: "b",
		MyEnumC: "c",
		MyEnumD: "d",
		MyEnumE: "e",
		MyEnumF: "f",
	}
	_MyEnumStringToValue = map[string]MyEnum{
		"a": MyEnumA,
		"b": MyEnumB,
		"c": MyEnumC,
		"d": MyEnumD,
		"e": MyEnumE,
		"f": MyEnumF,
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

var (
	_YourEnumValueToString = map[YourEnum]string{
		YourEnumA: "a",
		YourEnumB: "b",
	}
	_YourEnumStringToValue = map[string]YourEnum{
		"a": YourEnumA,
		"b": YourEnumB,
	}
)

func (e YourEnum) String() string {
	if str, ok := _YourEnumValueToString[e]; ok {
		return str
	}
	return fmt.Sprintf("YourEnum(%d)", e)
}

func (e YourEnum) IsVaild() bool {
	_, ok := _YourEnumValueToString[e]
	return ok
}

func ParseYourEnum(s string) (YourEnum, error) {
	if val, ok := _YourEnumStringToValue[s]; ok {
		return val, nil
	}
	return YourEnum(0), fmt.Errorf("invalid value %q for YourEnum", s)
}

func ParseYourEnumOr(s string, or YourEnum) YourEnum {
	val, err := ParseYourEnum(s)
	if err != nil {
		return or
	}
	return val
}

func MustParseYourEnum(s string) YourEnum {
	val, err := ParseYourEnum(s)
	if err != nil {
		panic(err)
	}
	return val
}

func (e YourEnum) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _YourEnumValueToString[e]; !ok {
		s = fmt.Sprintf("YourEnum(%d)", e)
	}
	return json.Marshal(s)
}

func (e *YourEnum) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := ParseYourEnum(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*e = newval
	return nil
}

func (e *YourEnum) Parse(s string) error {
	if val, ok := _YourEnumStringToValue[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid value %q for YourEnum", s)
}
