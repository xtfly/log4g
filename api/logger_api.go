package api

import (
	"context"
	"fmt"
	"time"
)

// -----------------------------
// --------Logger API-----------
// -----------------------------

// Writer for a logger
type Writer interface {
	// Tracef formats message according to format specifier
	// and writes to log with level = Trace.
	Tracef(fmt string, args ...interface{})

	// Trace formats message using the default formats for its operands
	// and writes to log with level = Trace
	Trace(msg ...interface{})

	// Debugf formats message according to format specifier
	// and writes to log with level = Debug.
	Debugf(fmt string, args ...interface{})

	// Debug formats message using the default formats for its operands
	// and writes to log with level = Debug
	Debug(msg ...interface{})

	// Infof formats message according to format specifier
	// and writes to log with level = Info.
	Infof(fmt string, args ...interface{})

	// Info formats message using the default formats for its operands
	// and writes to log with level = Info
	Info(msg ...interface{})

	// Warnf formats message according to format specifier
	// and writes to log with level = Warn.
	Warnf(fmt string, args ...interface{})

	// Warn formats message using the default formats for its operands
	// and writes to log with level = Warn
	Warn(msg ...interface{})

	// Errorf formats message according to format specifier
	// and writes to log with level = Error.
	Errorf(fmt string, args ...interface{})

	// Error formats message using the default formats for its operands
	// and writes to log with level = Error
	Error(msg ...interface{})

	// Criticalf formats message according to format specifier
	// and writes to log with level = Critical.
	Criticalf(fmt string, args ...interface{})

	// Critical formats message using the default formats for its operands
	// and writes to log with level = Critical
	Critical(msg ...interface{})

	// Printf message to logger using specified level
	Printf(lvl Level, fmt string, args ...interface{})
}

// Field is logging message extend field
type Field struct {
	Key   string
	Value interface{}
}

// Logger represents struct capable of logging messages
type Logger interface {
	// This method is used to add some extension fields before invoking the Tracef/Debugf method
	WithFields(fields ...Field) Writer
	WithCtx(ctx context.Context) Writer

	TraceEnabled() bool
	DebugEnabled() bool
	InfoEnabled() bool
	WarnEnabled() bool
	ErrorEnabled() bool
	CriticalEnabled() bool
	LevelEnabled(lvl Level) bool

	// SetOutputs ...
	SetOutputs(outputs []Output)

	// SetLevel to set the level of the logger
	SetLevel(lvl Level)

	// Level return the level of the logger
	Level() Level

	// SetCallerSkip to set the caller stack level will being skip
	SetCallerSkip(skip int)

	Writer
}

// Factory create or retrieve a logger by given name
type Factory interface {
	GetLogger(name string) Logger
}

// Event which is created when you logging and will send to Output,
// The Output implementer uses the Formatter to format the event to a logging content.
type Event struct {
	Format    string
	Arguments []interface{}
	Name      string
	Level     Level
	Time      time.Time
	CallDepth int
	Ctx       context.Context
}

// Message return a string which format by param 'Format' and 'Arguments'
func (e *Event) Message() string {
	msg := e.Format
	if len(e.Arguments) != 0 {
		if len(e.Format) != 0 {
			msg = fmt.Sprintf(e.Format, e.Arguments...)
		} else {
			msg = fmt.Sprint(e.Arguments...)
		}
	}
	return msg
}

// -----------------------------
// ------Formatter API----------
// -----------------------------

// Formatter formats events to logging content
type Formatter interface {
	// Format logging event to bytes
	Format(e *Event) []byte

	// CallerInfoFlag return caller info flag:
	// 0: none
	// 1: line, file
	// 2: func and above 1
	CallerInfoFlag() int
}

// -----------------------------
// ---------Output API----------
// -----------------------------

// Output appends contents to a Writer.
type Output interface {
	// Send a event to the output
	Send(e *Event)

	// SetFormatter set a formatter to the output
	SetFormatter(f Formatter)

	// CallerInfoFlag return the caller info flag by formatter
	CallerInfoFlag() int

	// Close the output and quit the loop routine
	Close()
}
