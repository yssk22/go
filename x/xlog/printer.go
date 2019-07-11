package xlog

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Printer is a temporary object to generate a log text within lazy formatting context.
type Printer struct {
	buff bytes.Buffer
}

// Printf to format printf style.
func (p *Printer) Printf(s string, v ...interface{}) {
	fmt.Fprintf(&p.buff, s, v...)
}

// Println to format println style.
func (p *Printer) Println(v ...interface{}) {
	fmt.Fprintln(&p.buff, v...)
}

type printerFunc func(*Printer)

func (f printerFunc) String() string {
	p := &Printer{}
	f(p)
	return p.buff.String()
}

func (f printerFunc) MarshalJSON() ([]byte, error) {
	p := &Printer{}
	f(p)
	return json.Marshal(p.buff.String())
}
