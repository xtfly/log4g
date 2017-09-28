package main

import (
	"fmt"

	"github.com/xtfly/log4g"
)

func main() {

	err := log.GetManager().LoadConfigFile("log4g.yaml")
	fmt.Errorf("%v", err)

	dlog := log.GetLogger("a/b")
	dlog.Debug("message")
	dlog.Info("info message")

	log.GetManager().Close()

}
