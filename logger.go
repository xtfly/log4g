package log

import (
	"context"
	"time"
)

const (
	callerSkip = 1
)

// defLogger 默认日志
type defLogger struct {
	name       string     // 日志名称
	level      Level      // 日志开启的级别
	parent     *defLogger // 日志的父一级
	outputs    []Output   // 日志的Appender列表
	callerSkip int        // caller skip depth

	*defWriter
}

func newLogger(name string) *defLogger {
	l := &defLogger{
		name:  name,
		level: Uninitialized,
	}
	w := &defWriter{logger: l, ctx: context.Background()}
	l.defWriter = w
	return l
}

func (l *defLogger) WithFields(fields ...Field) Writer {
	ctx := context.Background()
	for i := range fields {
		f := &fields[i]
		ctx = context.WithValue(ctx, f.Key, f.Value)
	}
	w := &defWriter{logger: l, ctx: ctx}
	return w
}

func (l *defLogger) WithCtx(ctx context.Context) Writer {
	return &defWriter{logger: l, ctx: ctx}
}

func (l *defLogger) TraceEnabled() bool {
	return l.LevelEnabled(Trace)
}

func (l *defLogger) DebugEnabled() bool {
	return l.LevelEnabled(Debug)
}

func (l *defLogger) InfoEnabled() bool {
	return l.LevelEnabled(Info)
}

func (l *defLogger) WarnEnabled() bool {
	return l.LevelEnabled(Warn)
}

func (l *defLogger) ErrorEnabled() bool {
	return l.LevelEnabled(Error)
}

func (l *defLogger) CriticalEnabled() bool {
	return l.LevelEnabled(Critical)
}

func (l *defLogger) LevelEnabled(lvl Level) bool {
	return lvl >= l.Level()
}

func (l *defLogger) SetLevel(lvl Level) {
	l.level = lvl
}

func (l *defLogger) Level() Level {
	cl := l
	for cl != nil {
		if cl.level != Uninitialized {
			return cl.level
		}
		cl = cl.parent
	}
	return Off
}

func (l *defLogger) SetCallerSkip(skip int) {
	l.callerSkip = callerSkip + skip
}

func (l *defLogger) SetOutputs(outputs []Output) {
	l.outputs = outputs
}

type defWriter struct {
	logger *defLogger
	ctx    context.Context
}

func (l *defWriter) Tracef(fmt string, args ...interface{}) {
	l.Printf(Trace, fmt, args...)
}

func (l *defWriter) Trace(msg string) {
	l.Printf(Trace, msg)
}

func (l *defWriter) Debugf(fmt string, args ...interface{}) {
	l.Printf(Debug, fmt, args...)
}

func (l *defWriter) Debug(msg string) {
	l.Printf(Debug, msg)
}

func (l *defWriter) Infof(fmt string, args ...interface{}) {
	l.Printf(Info, fmt, args...)
}

func (l *defWriter) Info(msg string) {
	l.Printf(Info, msg)
}

func (l *defWriter) Warnf(fmt string, args ...interface{}) {
	l.Printf(Warn, fmt, args...)
}

func (l *defWriter) Warn(msg string) {
	l.Printf(Warn, msg)
}

func (l *defWriter) Errorf(fmt string, args ...interface{}) {
	l.Printf(Error, fmt, args...)
}

func (l *defWriter) Error(msg string) {
	l.Printf(Error, msg)
}

func (l *defWriter) Criticalf(fmt string, args ...interface{}) {
	l.Printf(Critical, fmt, args...)
}

func (l *defWriter) Critical(msg string) {
	l.Printf(Critical, msg)
}

func (l *defWriter) Printf(lvl Level, fmt string, args ...interface{}) {
	l.write(l.logger.name, l.logger.callerSkip, lvl, fmt, args...)
}

func (l *defWriter) write(name string, skip int, lvl Level, fmt string, args ...interface{}) {
	if !l.logger.LevelEnabled(lvl) {
		return
	}

	if l.logger.outputs == nil || len(l.logger.outputs) == 0 {
		if l.logger.parent != nil {
			l.logger.parent.write(name, skip+1, lvl, fmt, args...)
			return
		}
	}

	// 生成事件
	evt := &Event{
		Time:      time.Now(),
		Name:      name,
		Level:     lvl,
		Format:    fmt,
		Arguments: args,
		CallDepth: skip,
		Ctx:       l.ctx,
	}

	// 分发给后端的Output
	for _, v := range l.logger.outputs {
		v.Send(evt)
	}
}
