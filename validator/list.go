package validator

// Func is a function to validate the value `v`
type Func func(v interface{}) *FieldError

// List is a list of Func
type List struct {
	funcs []Func
}

// NewList returns a new *List
func NewList() *List {
	return &List{make([]Func, 0)}
}

// Required is a function to validate the object is not a nil
func (l *List) Required() *List {
	l.funcs = append(l.funcs, requiredFunc)
	return l
}

// Min validates the value is more than or equal to `n`
func (l *List) Min(n int64) *List {
	l.funcs = append(l.funcs, minFunc(n))
	return l
}

// Max validates the value is more than or equal to `n`
func (l *List) Max(n int64) *List {
	l.funcs = append(l.funcs, maxFunc(n))
	return l
}

// Match validates the value matches `str`
func (l *List) Match(str string) *List {
	l.funcs = append(l.funcs, matchFunc(str))
	return l
}

// Unmatch validates the value does not matche `str`
func (l *List) Unmatch(str string) *List {
	l.funcs = append(l.funcs, unmatchFunc(str))
	return l
}

// Func to run the given custom validation fucntion `f`.
func (l *List) Func(f func(v interface{}) *FieldError) *List {
	l.funcs = append(l.funcs, f)
	return l
}
