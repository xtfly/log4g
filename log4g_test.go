package log

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func getOutput(id string) *memoryOutput {
	dm := gmanager.(*defManager)
	return dm.outputs[id].(*memoryOutput)
}

func init() {
	cfg := &Config{
		Loggers: []CfgLogger{
			{Name: "test", Level: "all", OutputNames: []string{"m1"}},
			{Name: "test2", Level: "debug", OutputNames: []string{"m2"}},
		},
		Formats: []CfgFormat{
			{Type: "format", Name: "f1", Properties: map[string]string{"layout": "%{module}|%{lvl}>>%{msg}"}},
			{Type: "format", Name: "f2", Properties: map[string]string{"layout": "%{f1:-5s}|%{f2:-5s}|%{module}|%{lvl:5s}>>%{msg}"}},
		},
		Outputs: []CfgOutput{
			{Type: "memory", Name: "m1", FormatName: "f1"},
			{Type: "memory", Name: "m2", FormatName: "f2"},
		},
	}

	gmanager.SetConfig(cfg)
}

func TestLoggerFormat(t *testing.T) {
	log := GetLogger("test")
	log.Debug("xxxx")

	mo := getOutput("m1")
	assert.Equal(t, mo.String(), "test|DBG>>xxxx")

	mo.buf.Truncate(0)
	log.Debugf("xxxx  %d", 1)
	assert.Equal(t, mo.String(), "test|DBG>>xxxx  1")

	mo.buf.Truncate(0)
	log.Trace("xxxx")
	assert.Equal(t, mo.String(), "test|TAC>>xxxx")

	mo.buf.Truncate(0)
	log.Tracef("xxxx  %d", 1)
	assert.Equal(t, mo.String(), "test|TAC>>xxxx  1")

	mo.buf.Truncate(0)
	log.Info("xxxx")
	assert.Equal(t, mo.String(), "test|INF>>xxxx")

	mo.buf.Truncate(0)
	log.Infof("xxxx  %d", 1)
	assert.Equal(t, mo.String(), "test|INF>>xxxx  1")

	mo.buf.Truncate(0)
	log.Warn("xxxx")
	assert.Equal(t, mo.String(), "test|WRN>>xxxx")

	mo.buf.Truncate(0)
	log.Warnf("xxxx  %d", 1)
	assert.Equal(t, mo.String(), "test|WRN>>xxxx  1")

	mo.buf.Truncate(0)
	log.Error("xxxx")
	assert.Equal(t, mo.String(), "test|ERR>>xxxx")

	mo.buf.Truncate(0)
	log.Errorf("xxxx  %d", 1)
	assert.Equal(t, mo.String(), "test|ERR>>xxxx  1")

	mo.buf.Truncate(0)
	log.Critical("xxxx")
	assert.Equal(t, mo.String(), "test|CRI>>xxxx")

	mo.buf.Truncate(0)
	log.Criticalf("xxxx  %d", 1)
	assert.Equal(t, mo.String(), "test|CRI>>xxxx  1")

	mo.buf.Truncate(0)
}

func TestLoggerFields(t *testing.T) {
	log := GetLogger("test2")
	w := log.WithFields(Field{"f1", "v1"}, Field{"f2", "v2"})
	w.Debug("xxxx")

	mo := getOutput("m2")
	assert.Equal(t, mo.String(), "v1   |v2   |test2|  DBG>>xxxx")
}


func TestLoggerParent(t *testing.T) {
	log := GetLogger("test")
	log.Debug("xxxx")

	mo := getOutput("m1")
	mo.buf.Truncate(0)

	log = GetLogger("test/1/2")
	log.Debug("xxxx")

	assert.Equal(t, mo.String(), "test/1/2|DBG>>xxxx")

	mo.buf.Truncate(0)
}