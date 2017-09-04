package log

import (
	"log"
	"sync"
)

const (
	// rootLoggerName is the name of the root logger.
	rootLoggerName    = "root"
	packageSepBySlash = '/'
	packageSepByDot   = '.'
)

// factory implements Factory interface.
type factory struct {
	sync.Mutex
	manager Manager
	root    *defLogger
	loggers map[string]*defLogger
}

func (f *factory) GetLogger(name string) Logger {
	f.Lock()
	defer f.Unlock()

	if name == "" || name == rootLoggerName {
		return f.getRootLogger()
	}

	l, ok := f.loggers[name]
	if !ok {
		l = f.createLogger(name, f.getParent(name))
		if ops, lvl, err := f.manager.GetLoggerOutputs(name); err != nil {
			log.Println("create logger error, ", err)
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
		f.root.SetLevel(Debug)
		console, _ := NewConsoleOutput(nil)
		f.root.SetOutputs([]Output{console})
	} else {
		f.root.SetLevel(lvl)
		f.root.SetOutputs(ops)
	}

	f.loggers[rootLoggerName] = f.root
	return f.root
}

// newFactory return a instance of Factory
func newFactory(manager Manager) Factory {
	factory := &factory{
		loggers: make(map[string]*defLogger),
		manager: manager,
	}

	return factory
}
