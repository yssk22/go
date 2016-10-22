package xlog

// Level is a enum for log Level
//go:generate enum -type=Level
type Level int

// Available Level values
const (
	LevelAll Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelNone
)

var minFilterFuncs = map[Level]FilterFunc{
	LevelAll:   minFilterFunc(LevelAll),
	LevelTrace: minFilterFunc(LevelTrace),
	LevelDebug: minFilterFunc(LevelDebug),
	LevelInfo:  minFilterFunc(LevelInfo),
	LevelWarn:  minFilterFunc(LevelWarn),
	LevelError: minFilterFunc(LevelError),
	LevelFatal: minFilterFunc(LevelFatal),
	LevelNone:  minFilterFunc(LevelNone),
}

func minFilterFunc(l Level) FilterFunc {
	return func(r *Record) bool {
		return r.Level >= l
	}
}

// Filter returns a Sink that filter records whose level is >= l
func (l Level) Filter(s Sink) Sink {
	return minFilterFuncs[l].Pipe(s)
}
