package internal

import (
	"context"
	"log"
	"time"

	"github.com/xtfly/log4g/api"
)

const (
	callerSkip = 3
)

// defLogger is default logger implements interface Logger
type defLogger struct {
	name           string       // 日志名称
	level          api.Level    // 日志开启的级别
	parent         *defLogger   // 日志的父一级
	outputs        []api.Output // 日志的Output列表
	callerSkip     int          // caller skip depth
	callerInfoFlag int          //

	*defWriter
}

func newLogger(name string) *defLogger {
	l := &defLogger{
		name:       name,
		level:      api.Uninitialized,
		callerSkip: callerSkip,
	}
	w := &defWriter{logger: l, ctx: context.Background()}
	l.defWriter = w
	return l
}

func (l *defLogger) WithFields(fields ...api.Field) api.Writer {
	ctx := context.Background()
	for i := range fields {
		f := &fields[i]
		ctx = context.WithValue(ctx, f.Key, f.Value)
	}
	w := &defWriter{logger: l, ctx: ctx}
	return w
}

func (l *defLogger) WithCtx(ctx context.Context) api.Writer {
	return &defWriter{logger: l, ctx: ctx}
}

func (l *defLogger) TraceEnabled() bool {
	return l.LevelEnabled(api.Trace)
}

func (l *defLogger) DebugEnabled() bool {
	return l.LevelEnabled(api.Debug)
}

func (l *defLogger) InfoEnabled() bool {
	return l.LevelEnabled(api.Info)
}

func (l *defLogger) WarnEnabled() bool {
	return l.LevelEnabled(api.Warn)
}

func (l *defLogger) ErrorEnabled() bool {
	return l.LevelEnabled(api.Error)
}

func (l *defLogger) CriticalEnabled() bool {
	return l.LevelEnabled(api.Critical)
}

func (l *defLogger) LevelEnabled(lvl api.Level) bool {
	return lvl >= l.Level()
}

func (l *defLogger) SetLevel(lvl api.Level) {
	l.level = lvl
}

func (l *defLogger) Level() api.Level {
	cl := l
	for cl != nil {
		if cl.level != api.Uninitialized {
			return cl.level
		}
		cl = cl.parent
	}
	return api.Off
}

func (l *defLogger) SetCallerSkip(skip int) {
	l.callerSkip = callerSkip + skip
}

func (l *defLogger) SetOutputs(outputs []api.Output) {
	l.outputs = outputs
	for _, op := range outputs {
		if l.callerInfoFlag < op.CallerInfoFlag() {
			l.callerInfoFlag = op.CallerInfoFlag()
		}
	}
}

type defWriter struct {
	logger *defLogger
	ctx    context.Context
}

func (l *defWriter) Tracef(fmt string, args ...interface{}) {
	l.Printf(api.Trace, fmt, args...)
}

func (l *defWriter) Trace(msg ...interface{}) {
	l.Printf(api.Trace, "", msg...)
}

func (l *defWriter) Debugf(fmt string, args ...interface{}) {
	l.Printf(api.Debug, fmt, args...)
}

func (l *defWriter) Debug(msg ...interface{}) {
	l.Printf(api.Debug, "", msg...)
}

func (l *defWriter) Infof(fmt string, args ...interface{}) {
	l.Printf(api.Info, fmt, args...)
}

func (l *defWriter) Info(msg ...interface{}) {
	l.Printf(api.Info, "", msg...)
}

func (l *defWriter) Warnf(fmt string, args ...interface{}) {
	l.Printf(api.Warn, fmt, args...)
}

func (l *defWriter) Warn(msg ...interface{}) {
	l.Printf(api.Warn, "", msg...)
}

func (l *defWriter) Errorf(fmt string, args ...interface{}) {
	l.Printf(api.Error, fmt, args...)
}

func (l *defWriter) Error(msg ...interface{}) {
	l.Printf(api.Error, "", msg...)
}

func (l *defWriter) Criticalf(fmt string, args ...interface{}) {
	l.Printf(api.Critical, fmt, args...)
}

func (l *defWriter) Critical(msg ...interface{}) {
	l.Printf(api.Critical, "", msg...)
}

func (l *defWriter) Printf(lvl api.Level, fmt string, args ...interface{}) {
	l.write(l.logger.name, l.logger.callerSkip, lvl, fmt, args...)
}

func (l *defWriter) write(name string, skip int, lvl api.Level, fmt string, args ...interface{}) {
	if !l.logger.LevelEnabled(lvl) {
		return
	}

	if len(l.logger.outputs) == 0 {
		if l.logger.parent != nil {
			l.logger.parent.write(name, skip+1, lvl, fmt, args...)
			return
		}
		log.Println("Warnning: not find outputs and parent for logger " + name)
	}

	// create a new logging event
	evt := &api.Event{
		Time:      time.Now(),
		Name:      name,
		Level:     lvl,
		Format:    fmt,
		Arguments: args,
		CallDepth: skip,
		Ctx:       l.ctx,
	}

	if l.logger.callerInfoFlag == ciFuncFlag {
		getCallInfo(evt, true)
	} else if l.logger.callerInfoFlag == ciFileFlag {
		getCallInfo(evt, false)
	}

	// dispatch event to all outputs
	for _, v := range l.logger.outputs {
		v.Send(evt)
	}
}
