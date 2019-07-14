package internal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/xtfly/log4g/api"
)

var (
	gmanager = newManager()
	gfactory = newFactory(gmanager)
)

func init() {
	gmanager.RegisterFormatterCreator(typeText, NewTextFormatter)

	gmanager.RegisterOutputCreator(typeConsole, NewConsoleOutput)
	gmanager.RegisterOutputCreator(typeMemory, NewMemoryOutput)
	gmanager.RegisterOutputCreator(typeRollingSize, NewRollingOutput)
	gmanager.RegisterOutputCreator(typeRollingTime, NewRollingOutput)
	gmanager.RegisterOutputCreator(typeSyslog, NewSyslogOutput)

	// default config
	cfg := &api.Config{
		Loggers: []api.CfgLogger{
			{Name: "root", Level: "all", OutputNames: []string{"c1"}},
		},
		Formats: []api.CfgFormat{
			{"type": "text", "name": "f1", "layout": "%{time}|%{level:-5s}|%{pid:6d}|%{module} >> %{msg} (%{longfunc}@%{longpkg}/%{shortfile}:%{line}) \n"},
		},
		Outputs: []api.CfgOutput{
			{"type": "console", "name": "c1", "format": "f1"},
		},
	}
	_ = gmanager.SetConfig(cfg)

	listenSig := make(chan os.Signal, 1)
	signal.Notify(listenSig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-listenSig
		gmanager.Close()
	}()
}

// GetLogger return the instance that implements Logger interface specified by name,
// name like a/b/c or a.b.c ,
// Logger named 'a/b' is the parent of Logger named 'a/b/c'
func GetLogger(name string) api.Logger {
	return gfactory.GetLogger(name)
}

// GetManager return the instance that implements Manager interface
func GetManager() api.Manager {
	return gmanager
}
