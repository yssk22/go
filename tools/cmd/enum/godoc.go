// Enum is a tool to automate the creation of methods for enum types.
//
// Setup:
//
//     go get github.com/yssk22/go/tools/cmd/enum/
//
// Usage: Define enum type using type alias and values by constants with the type name prefix.
//
//     // my_enum.go
//
//     //go:generate enum -type=MyEnum
//     type MyEnum int
//
//     const (
// 	       MyEnumReady MyEnum = iota
// 	       MyEnumRunning
// 	       MyEnumSuccess
// 	       MyEnumFailure
//     )
//
// Then you can generate 4 methos by `go generate` command` to use MyEnum as enum.
//
//     - func ParseMyEnum(s string) (MyEnum, error)
//     - func ParseMyEnumOr(s string, e MyEnum) (MyEnum)
//     - func (MyEnum) String()
//     - func (MyEnum) MarshalJSON() ([]byte, error)
//     - func (*MyEnum) UnmarshalJSON() ([]byte, error)
//
package main
