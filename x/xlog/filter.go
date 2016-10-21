package xlog

// LevelFilter returns a sink to pipeline for filter logs by Level to the next sink
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

// NameFilter returns a sink to pipeline for filter logs by Logger name to the next sink
func NameFilter(cfg map[string]Level, s Sink) Sink {
	return &nameFilter{
		cfg:  cfg,
		next: s,
	}
}

type nameFilter struct {
	cfg  map[string]Level
	next Sink
}

func (f *nameFilter) Write(r *Record) error {
	if level, ok := f.cfg[r.Name]; ok {
		if r.Level < level {
			return nil
		}
	}
	// discard
	return f.next.Write(r)
}
