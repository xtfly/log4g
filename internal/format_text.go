package internal

import (
	"bytes"

	"github.com/xtfly/log4g/api"
)

const (
	typeText = "text"
)

type textFormatter struct {
	layout       string
	strFormatter *StringFormatter
}

// NewTextFormatter return a Formatter instance using layout formatter string
func NewTextFormatter(cfg api.CfgFormat) (df api.Formatter, err error) {
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

// Format a logger event to byte array
func (f *textFormatter) Format(e *api.Event) []byte {
	var buf bytes.Buffer
	f.strFormatter.Format(e, &buf)
	return buf.Bytes()
}

// CallerInfoFlag return the max caller flag index
func (f *textFormatter) CallerInfoFlag() int {
	ret := 0
	for _, v := range f.strFormatter.parts {
		if r, ok := formatCallerFlags[v.verbName]; ok && r > ret {
			ret = r
		}
	}
	return ret
}
