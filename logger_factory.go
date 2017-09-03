package log

import (
	"log"
	"sync"
)

const (
	// rootLoggerName is the name of the root logger.
	rootLoggerName    = "root"
	packageSeparator  = '/'
	packageSeparator2 = '.'
)

// factory implements Factory interface.
type factory struct {
	sync.Mutex
	manager Manager
	root    *defLogger
	loggers map[string]*defLogger
}

func (f *factory) GetLogger(name string) Logger {
	if name == "" {
		return f.root
	}

	f.Lock()
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
	f.Unlock()

	return l
}

// getParent returns parent logger for given logger.
func (f *factory) getParent(name string) *defLogger {
	parent := f.root
	for i, c := range name {
		// Search for package separator character
		if c == packageSeparator || c == packageSeparator2 {
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

// newFactory return a instance of Factory
func newFactory(manager Manager) Factory {
	factory := &factory{
		root:    newLogger(rootLoggerName),
		loggers: make(map[string]*defLogger),
		manager: manager,
	}

	if ops, lvl, err := factory.manager.GetLoggerOutputs(rootLoggerName); err != nil {
		factory.root.SetLevel(Debug)
		console, _ := NewConsoleOutput(nil)
		factory.root.SetOutputs([]Output{console})
	} else {
		factory.root.SetLevel(lvl)
		factory.root.SetOutputs(ops)
	}

	factory.loggers[rootLoggerName] = factory.root
	return factory
}
