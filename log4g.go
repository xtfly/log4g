package log4g

import (
	"github.com/xtfly/log4g/api"
	"github.com/xtfly/log4g/internal"
)

// GetLogger return the instance that implements Logger interface specified by name,
// name like a/b/c or a.b.c ,
// Logger named 'a/b' is the parent of Logger named 'a/b/c'
func GetLogger(name string) api.Logger {
	return internal.GetLogger(name)
}

// GetManager return the instance that implements Manager interface
func GetManager() api.Manager {
	return internal.GetManager()
}
