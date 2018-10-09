package enum

// Spec represents Enum specification
type Spec struct {
	EnumName string
	Values   []Value
}

// Value represents a enum value
type Value struct {
	Name     string // constant name
	Value    int64  // enum value
	StrValue string // enum string value
}
