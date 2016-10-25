package xlog

import "io"

// Sink is an interface for log destination
type Sink interface {
	Write(*Record) error // Write a record to the destination
}

// NullSink is a sink to write nothing.
var NullSink = &nullSink{}

type nullSink struct{}

func (*nullSink) Write(*Record) error {
	return nil
}

// IOSink is an implementation of Sink that write log to io.Writer.
type IOSink struct {
	writer    io.Writer
	formatter Formatter
}

// NewIOSink returns a IOSink for w.
func NewIOSink(w io.Writer) *IOSink {
	return NewIOSinkWithFormatter(w, defaultIOFormatter)
}

// NewIOSinkWithFormatter returns a IOSink for w with f Formatter
func NewIOSinkWithFormatter(w io.Writer, f Formatter) *IOSink {
	if f == nil {
		f = defaultIOFormatter
	}
	return &IOSink{
		writer: w, formatter: f,
	}
}

// Write implements Sink#Write
func (s *IOSink) Write(r *Record) error {
	buff, err := s.formatter.Format(r)
	if err != nil {
		return err
	}
	_, err = s.writer.Write(buff)
	return err
}
