package internal

import (
	"log/syslog"

	"github.com/xtfly/log4g/api"
)

const (
	typeSyslog = "syslog"
)

type syslogOutput struct {
	w *syslog.Writer
	f api.Formatter
	t api.Level //threshold
}

func (o *syslogOutput) Send(e *api.Event) {
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
	case api.Trace, api.Debug:
		o.w.Debug(m)
	case api.Info:
		o.w.Info(m)
	case api.Warn:
		o.w.Warning(m)
	case api.Error:
		o.w.Err(m)
	case api.Critical:
		o.w.Crit(m)
	}
}

// SetFormatter ...
func (o *syslogOutput) SetFormatter(f api.Formatter) {
	o.f = f
}

// CallerInfoFlag return the formater max caller flag index
func (o *syslogOutput) CallerInfoFlag() int {
	if o.f != nil {
		return o.f.CallerInfoFlag()
	}
	return ciNoneFlog
}

// Close the syslog writer related to this output
func (o *syslogOutput) Close() {
	if o.w != nil {
		o.w.Close()
	}
}

// NewSyslogOutput return a output instance that output message to syslog
func NewSyslogOutput(cfg api.CfgOutput) (api.Output, error) {
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
