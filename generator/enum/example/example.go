package example

// MyEnum is an example of enum
// @enum
type MyEnum int

// @enum
type YourEnum int

const (
	MyEnumA MyEnum = iota
	MyEnumB
)

const MyEnumC, MyEnumD MyEnum = 11, 12

const MyEnumE, MyEnumF = MyEnum(13), MyEnum(14)

var NotMyEnumX, NotMyEnumY = MyEnum(100), MyEnum(123)

const NotMyEnumZ = 1

const (
	YourEnumA YourEnum = iota
	YourEnumB
)
