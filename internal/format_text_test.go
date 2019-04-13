package internal

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtfly/log4g/api"
)

func TestFormat(t *testing.T) {
	f, _ := NewTextFormatter(api.CfgFormat{"type": "text", "name": "f1",
		"layout": "%{shortfile} %{time:2006-01-02T15:04:05} %{level:.1s} %{module} %{msg}"})

	fbs := f.Format(&api.Event{
		Format:    "hello",
		Name:      "module",
		Level:     api.Debug,
		Ctx:       context.Background(),
		CallDepth: 2,
	})

	assert.Equal(t, string(fbs), "format_text.go 0001-01-01T00:00:00 D module hello")
}

func TestFormat2(t *testing.T) {
	var buf bytes.Buffer
	op := NewBaseOutput(&buf, api.All)
	f, _ := NewTextFormatter(api.CfgFormat{"type": "text", "name": "f1",
		"layout": "%{shortfile}"})
	op.SetFormatter(f)
	log := GetLogger("test_format")
	log.SetOutputs([]api.Output{op})

	log.Debug("xx")

	assert.Equal(t, buf.String(), "format_text_test.go")
}

func realFunc(log api.Logger) {
	log.Debug("xx")
}

type structFunc struct{}

func (structFunc) Log(log api.Logger) {
	log.Debug("xx")
}

func TestFormat3(t *testing.T) {
	var buf bytes.Buffer
	op := NewBaseOutput(&buf, api.All)
	f, _ := NewTextFormatter(api.CfgFormat{"type": "text", "name": "f1",
		"layout": "%{shortfunc}"})
	op.SetFormatter(f)
	log := GetLogger("test_format")
	log.SetOutputs([]api.Output{op})

	realFunc(log)

	assert.Equal(t, buf.String(), "realFunc")
}

func TestFormat4(t *testing.T) {
	var buf bytes.Buffer
	op := NewBaseOutput(&buf, api.All)
	f, _ := NewTextFormatter(api.CfgFormat{"type": "text", "name": "f1",
		"layout": "%{longfunc}"})
	op.SetFormatter(f)
	log := GetLogger("test_format")
	log.SetOutputs([]api.Output{op})

	structFunc{}.Log(log)

	assert.Equal(t, buf.String(), "structFunc.Log")
}

func TestFormat5(t *testing.T) {
	var buf bytes.Buffer
	op := NewBaseOutput(&buf, api.All)
	f, _ := NewTextFormatter(api.CfgFormat{"type": "text", "name": "f1",
		"layout": "%{shortfunc}"})
	op.SetFormatter(f)
	log := GetLogger("test_format")
	log.SetOutputs([]api.Output{op})

	var varFunc = func() string {
		log.Debug("xxx")
		return ""
	}
	varFunc()

	assert.Equal(t, buf.String(), "func1")
}
