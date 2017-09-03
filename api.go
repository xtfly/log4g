package log

import (
	"context"
	"time"
)

// 日志级别
type Level int

// 统一日志级别定义
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

type Writer interface {
	// 打印Trace级别的日志
	Tracef(fmt string, args ...interface{})
	Trace(msg string)

	// 打印Debug级别的日志
	Debugf(fmt string, args ...interface{})
	Debug(msg string)

	// 打印Info级别的日志
	Infof(fmt string, args ...interface{})
	Info(msg string)

	// 打印Warn级别的日志
	Warnf(fmt string, args ...interface{})
	Warn(msg string)

	// 打印Error级别的日志
	Errorf(fmt string, args ...interface{})
	Error(msg string)

	// 打印Critical级别的日志
	Criticalf(fmt string, args ...interface{})
	Critical(msg string)

	// Printf 打印指定级别的日志
	Printf(lvl Level, fmt string, args ...interface{})
}

// Field ...
type Field struct {
	Key   string
	Value interface{}
}

// Logger 普通日志接口
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

// Factory 获取日志的工厂接口
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

// Formatter 对日志事件格式化输出
type Formatter interface {
	// 格式化日志
	Format(e *Event) []byte
}

// FormatterFuncCreator ...
type FormatterFuncCreator func(arg *CfgFormat) (Formatter, error)

// Output appends contents to a Writer.
type Output interface {
	// Send 发送日志事件
	Send(e *Event)

	// SetFormatter ...
	SetFormatter(f Formatter)
}

// OutputFuncCreator ...
type OutputFuncCreator func(arg *CfgOutput) (Output, error)

type Manager interface {
	RegisterFormatterCreator(stype string, f FormatterFuncCreator)

	RegisterOutputCreator(stype string, o OutputFuncCreator)

	GetLoggerOutputs(name string) (ops []Output, lvl Level, err error)

	LoadConfig(file string) error

	SetConfig(cfg *Config) error
}

// Config ...
type Config struct {
	Formats []CfgFormat `yaml:"formats" json:"formats"`
	Outputs []CfgOutput `yaml:"output" json:"output"`
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
