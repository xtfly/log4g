package api

import (
	"context"
	"fmt"
	"time"
)

// -----------------------------
// --------Logger API-----------
// -----------------------------

// Level type for a logger
type Level int

// Log levels
const (
	Uninitialized Level = iota
	All
	Trace
	Debug
	Info
	Warn
	Error
	Critical
	Off
)

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

// Field is logger message extension field
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

	SetCallerSkip(skip int)

	Writer
}

// Factory represents ...
type Factory interface {
	GetLogger(name string) Logger
}

// Event is logging event which is created by logger and will send to Output,
// The Output implementer calls the Formatter interface to format the log output based on the event.
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

// Formatter represents ...
type Formatter interface {
	// Format logging event to bytes
	Format(e *Event) []byte

	// CallerInfoFlag return caller info flag:
	// 0: none
	// 1: line, file
	// 2: func and above 1
	CallerInfoFlag() int
}

// FormatterFuncCreator ...
type FormatterFuncCreator func(cfg CfgFormat) (Formatter, error)

// -----------------------------
// ---------Output API----------
// -----------------------------

// Output appends contents to a Writer.
type Output interface {
	// Send a event to the output
	Send(e *Event)

	// SetFormatter set a formatter to the output
	SetFormatter(f Formatter)

	// CallerInfoFlag ...
	CallerInfoFlag() int

	// Close the output and quit the loop routine
	Close()
}

// OutputFuncCreator ...
type OutputFuncCreator func(cfg CfgOutput) (Output, error)

// -----------------------------
// ---------Manager API---------
// -----------------------------

// Manager ...
type Manager interface {
	// RegisterFormatterCreator ..
	RegisterFormatterCreator(stype string, f FormatterFuncCreator)

	// RegisterOutputCreator ..
	RegisterOutputCreator(stype string, o OutputFuncCreator)

	// GetLoggerOutputs ..
	GetLoggerOutputs(name string) (ops []Output, lvl Level, err error)

	// LoadConfigFile ..
	LoadConfigFile(file string) error

	// LoadConfig ..
	LoadConfig(bs []byte, ext string) error

	// SetConfig ..
	SetConfig(cfg *Config) error

	// Close all output and wait all event write to outputs.
	Close()
}

// -----------------------------
// ---------Config API----------
// -----------------------------

// Config ...
type Config struct {
	Formats []CfgFormat `yaml:"formats" json:"formats"`
	Outputs []CfgOutput `yaml:"outputs" json:"outputs"`
	Loggers []CfgLogger `yaml:"loggers" json:"loggers"`
}

// GetCfgLogger return the point of CfgLogger which matched by name
func (c *Config) GetCfgLogger(name string) *CfgLogger {
	for _, l := range c.Loggers {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

// GetCfgOutput return the point of CfgOutput which matched by name
func (c *Config) GetCfgOutput(name string) CfgOutput {
	for _, l := range c.Outputs {
		if l.Name() == name {
			return l
		}
	}
	return nil
}

// GetCfgFormat return the point of CfgFormat which matched by name
func (c *Config) GetCfgFormat(name string) CfgFormat {
	for _, l := range c.Formats {
		if l.Name() == name {
			return l
		}
	}
	return nil
}

// CfgLogger ...
type CfgLogger struct {
	Name        string   `yaml:"name" json:"name"`
	Level       string   `yaml:"level" json:"level"`
	OutputNames []string `yaml:"outputs" json:"outputs"`
}

// CfgOutput ...
type CfgOutput map[string]string

// Name return the name of Output
func (c CfgOutput) Name() string {
	return c["name"]
}

// Type return the type of Output
func (c CfgOutput) Type() string {
	return c["type"]
}

// FormatName return the format name
func (c CfgOutput) FormatName() string {
	return c["format"]
}

// CfgFormat ...
type CfgFormat map[string]string

// Name return the name of Format
func (c CfgFormat) Name() string {
	return c["name"]
}

// Type return the type of Format
func (c CfgFormat) Type() string {
	return c["type"]
}
