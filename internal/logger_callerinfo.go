package internal

import (
	"context"
	"runtime"
	"strings"

	"github.com/xtfly/log4g/api"
)

const (
	ciNoneFlog = iota
	ciFileFlag
	ciFuncFlag
)

type callerInfo struct {
	file string
	line int
	pkg  string
	fun  string
	pc   uintptr
}

func getCallInfo(evt *api.Event, needFun bool) *callerInfo {
	var ci *callerInfo
	v := evt.Ctx.Value("__caller_info")
	if v != nil {
		ci = v.(*callerInfo)
	}

	if ci == nil {
		ci = &callerInfo{}
		var ok bool
		ci.pc, ci.file, ci.line, ok = runtime.Caller(evt.CallDepth + 1)
		if !ok {
			ci.file, ci.line = "???", 0
		}
		evt.Ctx = context.WithValue(evt.Ctx, "__caller_info", ci)
	}

	if needFun && ci.pkg == "" {
		if f := runtime.FuncForPC(ci.pc); f != nil {
			fs := f.Name()
			i := strings.LastIndex(fs, "/")
			j := strings.Index(fs[i+1:], ".")
			if j < 1 {
				ci.pkg, ci.fun = "???", "???"
			} else {
				ci.pkg, ci.fun = fs[:i+j+1], fs[i+j+2:]
			}
		}
	}

	return ci
}
