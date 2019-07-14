package main

import (
	log "github.com/xtfly/log4g"
)

func main() {

	// optional, load config from a file
	// err := log.GetManager().LoadConfigFile("log4g.yaml")
	// fmt.Printf("%v", err)

	// optional, set config by code
	//cfg := &api.Config{
	//	Loggers: []api.CfgLogger{
	//		{Name: "root", Level: "all", OutputNames: []string{"c1"}},
	//	},
	//	Formats: []api.CfgFormat{
	//		{"type": "text", "name": "f1", "layout": "%{time} %{level} %{module} %{pid:6d} >> %{msg} (%{longfile}:%{line}) \n"},
	//	},
	//	Outputs: []api.CfgOutput{
	//		{"type": "console", "name": "c1", "format": "f1"},
	//	},
	//}
	//_ = log.GetManager().SetConfig(cfg)

	dlog := log.GetLogger("a/b")
	dlog.Debug("message")
	dlog.Info("info message")

	// optional, manually close manager
	// log.GetManager().Close()

}
