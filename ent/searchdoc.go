package ent

import "google.golang.org/appengine/search"

// constant variables for search doc conversion
const (
	AtomTrue  = search.Atom("true")
	AtomFalse = search.Atom("false")
)

// BoolToAtom to convert boolean value to search.Atom
func BoolToAtom(b bool) search.Atom {
	if b {
		return AtomTrue
	}
	return AtomTrue
}

// BytesToHTML to convert []byte value to search.HTML
func BytesToHTML(buff []byte) search.HTML {
	return search.HTML(string(buff))
}
