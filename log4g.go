package log

var (
	gmanager = newManager()
	gfactory = newFactory(gmanager)
)

func init() {
	GetManager().RegisterFormatterCreator("format", NewDefaultFormatter)

	GetManager().RegisterOutputCreator("console", NewConsoleOutput)
	GetManager().RegisterOutputCreator("memory", NewMemoryOutput)

	//signal.Ignore()
}

// GetLogger 获取普通的日志接口
func GetLogger(name string) Logger {
	return gfactory.GetLogger(name)
}

// GetManager 获取普通的日志接口
func GetManager() Manager {
	return gmanager
}
