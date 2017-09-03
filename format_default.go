package log

import "bytes"

type defFormatter struct {
	layout       string
	strFormatter *StringFormatter
}

// NewDefaultFormatter return a Formatter instance using layout formatter string
func NewDefaultFormatter(arg *CfgFormat) (df Formatter, err error) {
	if arg == nil {
		panic("not set format config argument.")
	}
	obj := &defFormatter{
		layout:       arg.Properties["layout"],
		strFormatter: &StringFormatter{},
	}
	err = obj.strFormatter.Parser(obj.layout)
	df = obj
	return
}

// Format ...
func (f *defFormatter) Format(e *Event) []byte {
	var buf bytes.Buffer
	f.strFormatter.Format(e, &buf)
	return buf.Bytes()
}
