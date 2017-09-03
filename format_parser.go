// Copyright 2013, Örjan Persson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// base on https://github.com/op/go-logging/blob/master/format.go
package log

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	fmtVerbName = iota
	fmtVerbStatic
	fmtVerbTime
	fmtVerbExtend
)

const (
	verbTime   = "time"
	verbExtend = "_extend"
)

var (
	pid     = os.Getpid()
	program = filepath.Base(os.Args[0])
)

// FormatFunc takes raw logging event or any object
// that can be evaluated as string.
type FormatFunc func(evt *Event, part *part) interface{}

type part struct {
	verbType int //
	verbName string
	fmtStr   string
	fmtFunc  FormatFunc
	layout   string
}

// StringFormatter contains a list of parts which explains how to build the
// formatted string passed on to the logging backend.
type StringFormatter struct {
	parts []*part
}

var formatRe = regexp.MustCompile(`%{([a-zA-Z0-9]+)(?::(.*?[^\\]))?}`)

var (
	formatFuncs = map[string]FormatFunc{
		"pid":       pidFormatFunc,
		"program":   programFormatFunc,
		"module":    moduleFormatFunc,
		"msg":       messageFormatFunc,
		"level":     levelFormatFunc,
		"lvl":       lvlFormatFunc,
		"line":      lineFormatFunc,
		"longfile":  longfileFormatFunc,
		"shortfile": shortfileFormatFunc,
		"longpkg":   longpkgFormatFunc,
		"shortpkg":  shortpkgFormatFunc,
		"longfunc":  longfuncFormatFunc,
		"shortfunc": shortfuncFormatFunc,
		verbTime:    timeFormatFunc,
		verbExtend:  extendFormatFunc,
	}
)

// Parser 'verbs' specified in the format string.
//
// The verbs:
//
// General:
//     %{id}        Sequence number for log message (uint64).
//     %{pid}       Process id (int)
//     %{time}      Time when log occurred (time.Time)
//     %{level}     Log level (Level)
//     %{module}    Module (string)
//     %{program}   Basename of os.Args[0] (string)
//     %{message}   Message (string)
//     %{longfile}  Full file name and line number: /a/b/c/d.go:23
//     %{shortfile} Final file name element and line number: d.go:23
//
// For normal types, the output can be customized by using the 'verbs' defined
// in the fmt package, eg. '%{id:04d}' to make the id output be '%04d' as the
// format string.
//
// For time.Time, use the same layout as time.Format to change the time format
// when output, eg "2006-01-02T15:04:05.999Z-07:00".
//
// There's also a couple of experimental 'verbs'. These are exposed to get
// feedback and needs a bit of tinkering. Hence, they might change in the
// future.
//
// Experimental:
//     %{longpkg}   Full package path, eg. github.com/go-logging
//     %{shortpkg}  Base package path, eg. go-logging
//     %{longfunc}  Full function name, eg. littleEndian.PutUint32
//     %{shortfunc} Base function name, eg. PutUint32
//     %{callpath}  Call function path, eg. main.a.b.c
func (f *StringFormatter) Parser(layout string) error {
	// Find the boundaries of all %{vars}
	matches := formatRe.FindAllStringSubmatchIndex(layout, -1)
	if matches == nil {
		return errors.New("logger: invalid log format layout: " + layout)
	}

	// Collect all variables and static text for the format
	prev := 0
	for _, m := range matches {
		start, end := m[0], m[1]
		if start > prev {
			f.parts = append(f.parts, &part{verbType: fmtVerbStatic, fmtStr: layout[prev:start]})
		}

		name := layout[m[2]:m[3]]
		part := &part{verbType: fmtVerbName, verbName: name}
		//println(name)
		if ffunc, ok := formatFuncs[name]; ok {
			part.fmtFunc = ffunc
		} else {
			part.verbType = fmtVerbExtend
			part.fmtFunc = formatFuncs[verbExtend]
		}

		if m[4] != -1 {
			if name == verbTime {
				part.verbType = fmtVerbTime
				part.layout = layout[m[4]:m[5]]
			} else {
				part.fmtStr = "%" + layout[m[4]:m[5]]
			}
		}

		f.parts = append(f.parts, part)
		prev = end
	}

	end := layout[prev:]
	if end != "" {
		f.parts = append(f.parts, &part{verbType: fmtVerbStatic, fmtStr: end})
	}
	return nil
}

// Format logging event to output
func (f *StringFormatter) Format(e *Event, output io.Writer) {
	for _, part := range f.parts {
		switch part.verbType {
		case fmtVerbStatic:
			output.Write([]byte(part.fmtStr))
		case fmtVerbTime:
			output.Write([]byte(part.fmtFunc(e, part).(string)))
		default:
			// improve performance by call buffer.Write directly when fmtStr is empty
			if part.fmtStr == "" {
				if s, ok := part.fmtFunc(e, part).(string); ok {
					output.Write([]byte(s))
				} else {
					fmt.Fprintf(output, "%v", part.fmtFunc(e, part))
				}
			} else {
				fmt.Fprintf(output, part.fmtStr, part.fmtFunc(e, part))
			}
		}
	}
}

// --------------------------------------------------------

// %{pid}       Process id (int)
func pidFormatFunc(_ *Event, _ *part) interface{} {
	return pid
}

// %{program}   Basename of os.Args[0] (string)
func programFormatFunc(_ *Event, _ *part) interface{} {
	return program
}

// %{module}    Module (string)
func moduleFormatFunc(evt *Event, _ *part) interface{} {
	return evt.Name
}

// message
func messageFormatFunc(evt *Event, _ *part) interface{} {
	msg := evt.Format
	if len(evt.Arguments) != 0 {
		msg = fmt.Sprintf(evt.Format, evt.Arguments...)
	}
	return msg
}

// level
func levelFormatFunc(evt *Event, _ *part) interface{} {
	return LevelString(evt.Level)
}

// lvl
func lvlFormatFunc(evt *Event, _ *part) interface{} {
	return LvlString(evt.Level)
}

type callerInfo struct {
	file string
	line int
	pkg  string
	fun  string
	pc   uintptr
}

// %{line}  line number: 30
func lineFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, false)
	return ci.line
}

// %{longfile}  Full file name: /a/b/c/d.go
func longfileFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, false)
	return ci.file
}

// %{shortfile} Final file name element: d.go
func shortfileFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, false)
	return filepath.Base(ci.file)
}

// %{longpkg}  Full package path, eg. github.com//xtfy/log4g
func longpkgFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, true)
	return ci.pkg
}

// %{shortpkg}  Base package path, eg. log4g
func shortpkgFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, true)
	return path.Base(ci.pkg)
}

// %{longfunc}  Full function name, eg. littleEndian.PutUint32
func longfuncFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, true)
	return ci.fun
}

// %{shortfunc} Base function name, eg. PutUint32
func shortfuncFormatFunc(evt *Event, _ *part) interface{} {
	ci := getCallInfo(evt, true)
	i := strings.LastIndex(ci.fun, ".")
	return ci.fun[i+1:]
}

func getCallInfo(evt *Event, needFun bool) *callerInfo {
	ci := evt.Ctx.Value("__caller_info").(*callerInfo)
	if ci == nil {
		ci = &callerInfo{}
		var ok bool
		ci.pc, ci.file, ci.line, ok = runtime.Caller(evt.CallDepth)
		if !ok {
			ci.file, ci.line = "???", 0
		}
		evt.Ctx = context.WithValue(evt.Ctx, "__caller_info", ci)
	}

	if needFun {
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

// %{time} Time when log occurred (time.Time)
func timeFormatFunc(evt *Event, part *part) interface{} {
	if part.layout == "" {
		return evt.Time.Format(defaultTimeLayout)
	}
	return evt.Time.Format(part.layout)
}

// %{xxx} get value from context by xxx
func extendFormatFunc(evt *Event, part *part) interface{} {
	return evt.Ctx.Value(part.verbName)
}
