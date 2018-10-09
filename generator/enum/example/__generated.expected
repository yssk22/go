package example

import (
	"encoding/json"
	"fmt"
)

var (
	_MyEnumValueToString = map[MyEnum]string{
		MyEnumA: "aa",
		MyEnumB: "bb",
		MyEnumC: "cc",
		MyEnumD: "dd",
		MyEnumE: "ee",
		MyEnumF: "ff",
	}
	_MyEnumStringToValue = map[string]MyEnum{
		"aa": MyEnumA,
		"bb": MyEnumB,
		"cc": MyEnumC,
		"dd": MyEnumD,
		"ee": MyEnumE,
		"ff": MyEnumF,
	}
)

func (e MyEnum) String() string {
	if str, ok := _MyEnumValueToString[i]; ok {
		return str
	}
	return fmt.Sprintf("MyEnum(%d)", e)
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

func (i MyEnum) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _MyEnumValueToString[i]; !ok {
		s = fmt.Sprintf("MyEnum(%d)", i)
	}
	return json.Marshal(s)
}

func (i *MyEnum) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := ParseMyEnum(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*i = newval
	return nil
}

var (
	_YourEnumValueToString = map[YourEnum]string{
		YourEnumA: "aa",
		YourEnumB: "bb",
	}
	_YourEnumStringToValue = map[string]YourEnum{
		"aa": YourEnumA,
		"bb": YourEnumB,
	}
)

func (e YourEnum) String() string {
	if str, ok := _YourEnumValueToString[i]; ok {
		return str
	}
	return fmt.Sprintf("YourEnum(%d)", e)
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

func (i YourEnum) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _YourEnumValueToString[i]; !ok {
		s = fmt.Sprintf("YourEnum(%d)", i)
	}
	return json.Marshal(s)
}

func (i *YourEnum) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("invalid JSON string")
	}
	newval, err := ParseYourEnum(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*i = newval
	return nil
}
