package internal

import (
	"sync"

	"github.com/xtfly/log4g/api"
)

const (
	rootLoggerName    = "root"
	packageSepBySlash = '/'
	packageSepByDot   = '.'
)

// factory implements Factory interface.
type factory struct {
	sync.Mutex
	manager api.Manager
	root    *defLogger
	loggers map[string]*defLogger
}

func (f *factory) GetLogger(name string) api.Logger {
	f.Lock()
	defer f.Unlock()

	if name == "" || name == rootLoggerName {
		return f.getRootLogger()
	}

	l, ok := f.loggers[name]
	if !ok {
		l = f.createLogger(name, f.getParent(name))
		if ops, lvl, err := f.manager.GetLoggerOutputs(name); err != nil {
			//log.Println("WARN: ", err)
		} else {
			l.SetLevel(lvl)
			l.SetOutputs(ops)
		}
	}

	return l
}

// getParent returns parent logger for given logger.
func (f *factory) getParent(name string) *defLogger {
	parent := f.getRootLogger()
	for i, c := range name {
		// Search for package separator character
		if c == packageSepBySlash || c == packageSepByDot {
			parentName := name[0:i]
			if parentName != "" {
				parent = f.createLogger(parentName, parent)
			}
		}
	}
	return parent
}

// createLogger creates a new logger if not exist.
func (f *factory) createLogger(name string, parent *defLogger) *defLogger {
	l, ok := f.loggers[name]
	if !ok {
		l = newLogger(name)
		l.parent = parent
		f.loggers[name] = l
	}
	return l
}

func (f *factory) getRootLogger() *defLogger {
	if f.root != nil {
		return f.root
	}

	f.root = newLogger(rootLoggerName)
	if ops, lvl, err := f.manager.GetLoggerOutputs(rootLoggerName); err != nil {
		f.root.SetLevel(api.Debug)
		console, _ := NewConsoleOutput(nil)
		f.root.SetOutputs([]api.Output{console})
	} else {
		f.root.SetLevel(lvl)
		f.root.SetOutputs(ops)
	}

	f.loggers[rootLoggerName] = f.root
	return f.root
}

func (f *factory) notify() {
	f.Lock()
	defer f.Unlock()
	for _, k := range f.loggers {
		if ops, lvl, err := f.manager.GetLoggerOutputs(k.name); err != nil {
			//log.Println("WARN: ", err)
		} else {
			k.SetLevel(lvl)
			k.SetOutputs(ops)
		}
	}
}

// newFactory return a instance of Factory
func newFactory(manager api.Manager) api.Factory {
	factory := &factory{
		loggers: make(map[string]*defLogger),
		manager: manager,
	}
	(manager.(*defManager)).addConfigNotify(factory)
	return factory
}
