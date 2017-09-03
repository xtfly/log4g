package example

import "github.com/xtfly/log4g"

func main() {

	_ := log.GetManager().LoadConfig("log4g.yaml")

	dlog := log.GetLogger("a/b")
	dlog.Debug("message")

}
