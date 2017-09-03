package log

import (
	"fmt"
	"io"
)

const (
	defaultTimeLayout = "2006-01-02T15:04:05.000Z07:00"
	defaultLayout     = "%-5s [%s] %s: %s\n"
)

type baseOutput struct {
	w io.Writer
	f Formatter
}

func (o *baseOutput) Send(e *Event) {
	if o.f != nil {
		o.w.Write([]byte(o.f.Format(e)))
		return
	}

	var msg string
	if len(e.Arguments) == 0 {
		msg = e.Format
	} else {
		msg = fmt.Sprintf(e.Format, e.Arguments...)
	}

	fmt.Fprintf(o.w, defaultLayout,
		LevelString(e.Level),
		e.Time.Format(defaultTimeLayout),
		e.Name,
		msg)
}

func (o *baseOutput) SetFormatter(f Formatter) {
	o.f = f
}
