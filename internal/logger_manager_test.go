package internal

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/xtfly/log4g/api"
)

const (
	testCfgStr = `
formats:
  - name: f1
    type: text
    layout: "%{time} %{level} ${module} %{pid:6d} >> %{msg} (%{longfile}:%{line}) \n"
outputs:
  - name: c1
    type: console
    format: f1
    threshold: info
  - name: r1
    type: size_rolling_file
    format: f1
    file: log/rf.log
    file_perm: 0640
    back_perm: 0550
    dir_perm: 0650
    size: 1M
    backups: 5
    threshold: info
  - name: r2
    type: time_rolling_file
    format: f1
    file: log/rf.log
    file_perm: 0640
    back_perm: 0550
    dir_perm: 0650
    pattern: 2006-01-02
    backups: 5
    threshold: info
  - name: s1
    type: syslog
    format: f1
    prefix: module
loggers:
  - name: root
    level: info
    outputs: ["r1", "c1"]
  - name: a/b
    level: error
    outputs: ["r2", "s1"]
`
)

func TestManager(t *testing.T) {
	m := newManager()
	m.RegisterFormatterCreator(typeText, NewTextFormatter)
	m.RegisterOutputCreator(typeConsole, NewConsoleOutput)
	m.RegisterOutputCreator(typeMemory, NewMemoryOutput)
	m.RegisterOutputCreator(typeRollingSize, NewRollingOutput)
	m.RegisterOutputCreator(typeRollingTime, NewRollingOutput)
	m.RegisterOutputCreator(typeSyslog, NewSyslogOutput)

	err := m.LoadConfig([]byte(testCfgStr), "yaml")
	fmt.Printf("%v\n", err)

	ops, lvl, err := m.GetLoggerOutputs("root")
	assert.NoError(t, err)
	assert.Equal(t, api.Info, lvl)
	assert.Equal(t, 2, len(ops))

	ops, lvl, err = m.GetLoggerOutputs("a/b")
	assert.NoError(t, err)
	assert.Equal(t, api.Error, lvl)
	assert.Equal(t, 2, len(ops))

	m.Close()
	ops, lvl, err = m.GetLoggerOutputs("a/b")
	assert.NoError(t, err)
	assert.Equal(t, api.Error, lvl)
	assert.Equal(t, 2, len(ops))
}
