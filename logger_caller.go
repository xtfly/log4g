package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var workingDir = "/"

func init() {
	wd, err := os.Getwd()
	if err == nil {
		workingDir = filepath.ToSlash(wd) + "/"
	}
}

func extractCallerInfo(skip int) (fullPath string, shortPath string, funcName string, line int, err error) {
	pc, fp, ln, ok := runtime.Caller(skip)
	if !ok {
		err = fmt.Errorf("error during runtime.Caller")
		return
	}
	line = ln
	fullPath = fp
	if strings.HasPrefix(fp, workingDir) {
		shortPath = fp[len(workingDir):]
	} else {
		shortPath = fp
	}
	funcName = runtime.FuncForPC(pc).Name()
	if strings.HasPrefix(funcName, workingDir) {
		funcName = funcName[len(workingDir):]
	}
	return
}
