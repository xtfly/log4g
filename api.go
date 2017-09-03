package log

import (
	"context"
	"time"
)

// Log level type
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

// Writer logger writer
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

	// Printf 打印指定级别的日志
	Printf(lvl Level, fmt string, args ...interface{})
}

// Field ...
type Field struct {
	Key   string
	Value interface{}
}

// Logger represents structs capable of logging messages
type Logger interface {
	// 增加日志扩展字段，可以采用链式调用，返回本身
	// 此方法用于在调用Tracef/Debugf方法之前调用来增加一些扩展字段
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

	// 设置日志级别
	SetLevel(lvl Level)
	Level() Level

	SetCallerSkip(skip int)

	Writer
}

// Factory represents ...
type Factory interface {
	GetLogger(name string) Logger
}

// Event 日志事件，发给Output，
// Output实现者，基于事件来调用Formatter接口格式化日志输出
type Event struct {
	Format    string
	Arguments []interface{}
	Name      string
	Level     Level
	Time      time.Time
	CallDepth int
	Ctx       context.Context
}

// Formatter represents ...
type Formatter interface {
	// 格式化日志
	Format(e *Event) []byte
}

// FormatterFuncCreator ...
type FormatterFuncCreator func(arg *CfgFormat) (Formatter, error)

// Output appends contents to a Writer.
type Output interface {
	// Send a event to the output
	Send(e *Event)

	// SetFormatter set a formatter to the output
	SetFormatter(f Formatter)

	// Close the output and quit the loop routine
	Close()
}

// OutputFuncCreator ...
type OutputFuncCreator func(arg *CfgOutput) (Output, error)

type Manager interface {
	// RegisterFormatterCreator ..
	RegisterFormatterCreator(stype string, f FormatterFuncCreator)

	// RegisterOutputCreator ..
	RegisterOutputCreator(stype string, o OutputFuncCreator)

	// GetLoggerOutputs ..
	GetLoggerOutputs(name string) (ops []Output, lvl Level, err error)

	// LoadConfig ..
	LoadConfig(file string) error

	// SetConfig ..
	SetConfig(cfg *Config) error

	// Close all output and wait all event write to outputs.
	Close()
}

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

// GetCfgOutput return the point of CfgOutput which matched by id
func (c *Config) GetCfgOutput(id string) *CfgOutput {
	for _, l := range c.Outputs {
		if l.Name == id {
			return &l
		}
	}
	return nil
}

// GetCfgFormat return the point of CfgFormat which matched by id
func (c *Config) GetCfgFormat(id string) *CfgFormat {
	for _, l := range c.Formats {
		if l.Name == id {
			return &l
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

// CfgFormat ...
type CfgFormat struct {
	Name       string            `yaml:"name" json:"name"`
	Type       string            `yaml:"type" json:"type"`
	Properties map[string]string `yaml:"properties" json:"properties"`
}

// CfgOutput ...
type CfgOutput struct {
	Name       string            `yaml:"name" json:"name"`
	Type       string            `yaml:"type" json:"type"`
	FormatName string            `yaml:"format" json:"format"` // format id
	Properties map[string]string `yaml:"properties" json:"properties"`
}
