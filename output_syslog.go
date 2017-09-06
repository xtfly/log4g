package log

import (
	"log/syslog"
)

const (
	typeSyslog = "syslog"
)

type syslogOutput struct {
	w *syslog.Writer
	f Formatter
	t Level //threshold
}

func (o *syslogOutput) Send(e *Event) {
	if e.Level < o.t {
		return
	}

	m := ""
	if o.f != nil {
		m = string(o.f.Format(e))
	} else {
		m = e.Message()
	}

	switch e.Level {
	case Trace, Debug:
		o.w.Debug(m)
	case Info:
		o.w.Info(m)
	case Warn:
		o.w.Warning(m)
	case Error:
		o.w.Err(m)
	case Critical:
		o.w.Crit(m)
	}
}

// SetFormatter ...
func (o *syslogOutput) SetFormatter(f Formatter) {
	o.f = f
}

// CallerInfoFlag ...
func (o *syslogOutput) CallerInfoFlag() int {
	if o.f != nil {
		return o.f.CallerInfoFlag()
	}
	return ciNoneFlog
}

// Close ...
func (o *syslogOutput) Close() {
	if o.w != nil {
		o.w.Close()
	}
}

// NewSyslogOutput return a output instance that it print message to syslog
func NewSyslogOutput(cfg CfgOutput) (Output, error) {
	w, err := syslog.New(syslog.LOG_CRIT, cfg["prefix"])
	if err != nil {
		return nil, err
	}
	r := &syslogOutput{
		w: w,
		t: GetThresholdLvl(cfg["threshold"]),
	}
	return r, nil
}
