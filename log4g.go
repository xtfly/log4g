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
	GetManager().RegisterFormatterCreator("format", NewDefaultFormatter)

	GetManager().RegisterOutputCreator("console", NewConsoleOutput)
	GetManager().RegisterOutputCreator("memory", NewMemoryOutput)
	GetManager().RegisterOutputCreator("size_rolling_file", NewRollingOutput)
	GetManager().RegisterOutputCreator("time_rolling_file", NewRollingOutput)

	listenSig := make(chan os.Signal, 1)
	signal.Notify(listenSig, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range listenSig {
			GetManager().Close()
		}
	}()
}

// GetLogger 获取普通的日志接口
func GetLogger(name string) Logger {
	return gfactory.GetLogger(name)
}

// GetManager 获取普通的日志接口
func GetManager() Manager {
	return gmanager
}
