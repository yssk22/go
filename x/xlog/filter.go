package xlog

// FilterFunc
type FilterFunc func(r *Record) bool

// Pipe returns a *Filter that writes to Sink using FilterFunc
func (f FilterFunc) Pipe(s Sink) *Filter {
	return &Filter{
		f:    f,
		next: s,
	}
}

// Filter is a Sink to write records by condition to the next sink.
type Filter struct {
	f    FilterFunc
	next Sink
}

// Write implements `Sink#Write`
func (f *Filter) Write(r *Record) error {
	if f.f(r) {
		return f.next.Write(r)
	}
	return nil
}

// LevelFilter is FilterFunc that filter records whose level is under the given level.
var LevelFilter = func(min Level) FilterFunc {
	return func(r *Record) bool {
		return r.Level >= min
	}
}

// KeyLevelFilter is like LevelFilter but more flexible by key matching.
// If key is not found in `keyLevels`, default level l is used as fallback.
var KeyLevelFilter = func(keyLevels map[interface{}]Level, l Level) FilterFunc {
	return func(r *Record) bool {
		if lvl, ok := keyLevels[r.LoggerKey]; ok {
			return r.Level >= lvl
		}
		return r.Level >= l
	}
}
