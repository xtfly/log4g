package log

import (
	"os"
)

type consoleOutput struct {
	Output
}

// NewConsoleOutput return a output instance that it print message to stdio
func NewConsoleOutput(arg *CfgOutput) (Output, error) {
	r := &consoleOutput{}
	if arg != nil && arg.Properties["async"] == "true" {
		r.Output = NewAynscOutput(os.Stdout,
			GetQueueSize(arg.Properties["queue_size"]), GetBatchNum(arg.Properties["batch_num"]))
	} else {
		r.Output = NewBaseOutput(os.Stdout)
	}

	return r, nil
}
