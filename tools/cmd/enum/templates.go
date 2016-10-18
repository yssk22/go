package main

var supportedFuncs = []string{
	`func (i %[1]s) String() string {
	if str, ok := _%[1]sValueToString[i]; ok {
		return str
	}
	return fmt.Sprintf("%[1]s(%%d)", i)
}
`,
	`func Parse%[1]s(s string) (%[1]s, error) {
	if val, ok := _%[1]sStringToValue[s]; ok {
		return val, nil
	}
	return %[1]s(0), fmt.Errorf("Invalid value %%q for %[1]s", s)
}
`,
	`func Parse%[1]sOr(s string, or %[1]s) %[1]s {
	val, err :=  Parse%[1]s(s)
	if err != nil {
		return or
	}
	return val
}
`,
	`func (i %[1]s) MarshalJSON() ([]byte, error) {
	var s string
	var ok bool
	if s, ok = _%[1]sValueToString[i]; !ok {
		s = fmt.Sprintf("%[1]s(%%d)", i)
	}
	return json.Marshal(s)
}
`,
	`func (i *%[1]s) UnmarshalJSON(b []byte) (error) {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return fmt.Errorf("Invalid string")
	}
	newval, err := Parse%[1]s(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*i = newval
	return nil
}
`,
}
