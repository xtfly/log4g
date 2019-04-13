package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xtfly/log4g/api"
)

func getOutput(id string) *memoryOutput {
	dm := gmanager.(*defManager)
	return dm.outputs[id].(*memoryOutput)
}

func init() {
	cfg := &api.Config{
		Loggers: []api.CfgLogger{
			{Name: "test", Level: "all", OutputNames: []string{"m1"}},
			{Name: "test2", Level: "debug", OutputNames: []string{"m2"}},
		},
		Formats: []api.CfgFormat{
			{"type": "text", "name": "f1", "layout": "%{module}|%{lvl}>>%{msg}"},
			{"type": "text", "name": "f2", "layout": "%{f1:-5s}|%{f2:-5s}|%{module}|%{lvl:5s}>>%{msg}"},
		},
		Outputs: []api.CfgOutput{
			{"type": "memory", "name": "m1", "format": "f1"},
			{"type": "memory", "name": "m2", "format": "f2"},
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
	w := log.WithFields(api.Field{Key: "f1", Value: "v1"}, api.Field{Key: "f2", Value: "v2"})
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
