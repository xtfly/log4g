package log

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	gmanager = newManager()
	gfactory = newFactory(gmanager)
)

func init() {
	GetManager().RegisterFormatterCreator(typeText, NewTextFormatter)

	GetManager().RegisterOutputCreator(typeConsole, NewConsoleOutput)
	GetManager().RegisterOutputCreator(typeMemory, NewMemoryOutput)
	GetManager().RegisterOutputCreator(typeRollingSize, NewRollingOutput)
	GetManager().RegisterOutputCreator(typeRollingTime, NewRollingOutput)
	GetManager().RegisterOutputCreator(typeSyslog, NewSyslogOutput)

	listenSig := make(chan os.Signal, 1)
	signal.Notify(listenSig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-listenSig
		GetManager().Close()
	}()
}

// GetLogger return the instance that implements Logger interface specified by name,
// name like a/b/c or a.b.c ,
// Logger named 'a/b' is the parent of Logger named 'a/b/c'
func GetLogger(name string) Logger {
	return gfactory.GetLogger(name)
}

// GetManager return the instance that implements Manager interface
func GetManager() Manager {
	return gmanager
}
