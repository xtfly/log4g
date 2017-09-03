package log

import (
	"os"
)

type consoleOutput struct {
	*baseOutput
}

// NewConsoleOutput return a output instance that it print message to stdio
func NewConsoleOutput(_ *CfgOutput) (Output, error) {
	r := &consoleOutput{}
	r.baseOutput = &baseOutput{
		w: os.Stdout,
	}
	return r, nil
}
