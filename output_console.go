package log

import (
	"os"
)

const (
	typeConsole = "console"
)

type consoleOutput struct {
	Output
}

// NewConsoleOutput return a output instance that it print message to stdio
func NewConsoleOutput(cfg CfgOutput) (Output, error) {
	r := &consoleOutput{}
	if cfg != nil && cfg["async"] == "true" {
		r.Output = NewAynscOutput(os.Stdout,
			GetQueueSize(cfg["queue_size"]), GetBatchNum(cfg["batch_num"]))
	} else {
		r.Output = NewBaseOutput(os.Stdout)
	}

	return r, nil
}
