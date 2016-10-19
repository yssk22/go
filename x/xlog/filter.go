package xlog

// LevelFilter returns a sink to pipeline for filter logs to the next sink
func LevelFilter(level Level, s Sink) Sink {
	return &levelFilter{
		level: level,
		next:  s,
	}
}

type levelFilter struct {
	level Level
	next  Sink
}

func (f *levelFilter) Write(r *Record) error {
	if r.Level < f.level {
		return nil
	}
	// discard
	return f.next.Write(r)
}
