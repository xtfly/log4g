package internal

import (
	"os"

	"github.com/xtfly/log4g/api"
)

const (
	typeConsole = "console"
)

type consoleOutput struct {
	api.Output
}

// NewConsoleOutput return a output instance that it print message to stdio
func NewConsoleOutput(cfg api.CfgOutput) (api.Output, error) {
	r := &consoleOutput{}
	if cfg != nil && cfg["async"] == "true" {
		r.Output = NewAsyncOutput(os.Stdout, GetThresholdLvl(cfg["threshold"]),
			GetQueueSize(cfg["queue_size"]), GetBatchNum(cfg["batch_num"]))
	} else {
		r.Output = NewBaseOutput(os.Stdout, GetThresholdLvl(cfg["threshold"]))
	}

	return r, nil
}
