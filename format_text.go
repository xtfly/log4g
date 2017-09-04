package log

import "bytes"

const (
	typeText = "text"
)

type textFormatter struct {
	layout       string
	strFormatter *StringFormatter
}

// NewTextFormatter return a Formatter instance using layout formatter string
func NewTextFormatter(cfg CfgFormat) (df Formatter, err error) {
	if cfg == nil {
		panic("not set format config argument.")
	}
	obj := &textFormatter{
		layout:       cfg["layout"],
		strFormatter: &StringFormatter{},
	}
	err = obj.strFormatter.Parser(obj.layout)
	df = obj
	return
}

// Format ...
func (f *textFormatter) Format(e *Event) []byte {
	var buf bytes.Buffer
	f.strFormatter.Format(e, &buf)
	return buf.Bytes()
}
