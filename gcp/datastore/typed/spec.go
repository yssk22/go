package typed

// Spec is a specificaiton for datastore entity
type Spec struct {
	StructName     string // struct name
	KindName       string // entity kind name (usually same as StructName) but different if as=XX is specfiied
	KeyField       string
	TimestampField string
	IsSearchable   bool
	Fields         []*FieldSpec
	QuerySpecs     []*QuerySpec
}

// FieldSpec is a specification for datasatore entity fields.
type FieldSpec struct {
	Name        string
	IsKey       bool
	IsID        bool
	IsTimestamp bool
	IsSearch    bool

	NoIndex bool
}

// QuerySpec is a specification for query
type QuerySpec struct {
	Name         string
	PropertyName string
	Type         string
}
