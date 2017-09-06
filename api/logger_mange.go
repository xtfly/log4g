package api

// -----------------------------
// ---------Manager API---------
// -----------------------------

// FormatterFuncCreator is function will to create a Formatter instance by configuration
type FormatterFuncCreator func(cfg CfgFormat) (Formatter, error)

// OutputFuncCreator is function will to create a Output instance by configuration
type OutputFuncCreator func(cfg CfgOutput) (Output, error)

// Manager is the configurations and creators holder
type Manager interface {
	// RegisterFormatterCreator ..
	RegisterFormatterCreator(stype string, f FormatterFuncCreator)

	// RegisterOutputCreator ..
	RegisterOutputCreator(stype string, o OutputFuncCreator)

	// GetLoggerOutputs ..
	GetLoggerOutputs(name string) (ops []Output, lvl Level, err error)

	// LoadConfigFile ..
	LoadConfigFile(file string) error

	// LoadConfig ..
	LoadConfig(bs []byte, ext string) error

	// SetConfig ..
	SetConfig(cfg *Config) error

	// Close all output and wait all event write to outputs.
	Close()
}

// -----------------------------
// ---------Config API----------
// -----------------------------

// Config struct aggregates all formatter, output and logger configurations
type Config struct {
	Formats []CfgFormat `yaml:"formats" json:"formats"`
	Outputs []CfgOutput `yaml:"outputs" json:"outputs"`
	Loggers []CfgLogger `yaml:"loggers" json:"loggers"`
}

// GetCfgLogger return the point of CfgLogger which matched by name
func (c *Config) GetCfgLogger(name string) *CfgLogger {
	for _, l := range c.Loggers {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

// GetCfgOutput return the point of CfgOutput which matched by name
func (c *Config) GetCfgOutput(name string) CfgOutput {
	for _, l := range c.Outputs {
		if l.Name() == name {
			return l
		}
	}
	return nil
}

// GetCfgFormat return the point of CfgFormat which matched by name
func (c *Config) GetCfgFormat(name string) CfgFormat {
	for _, l := range c.Formats {
		if l.Name() == name {
			return l
		}
	}
	return nil
}

// CfgLogger represents the configuration of a logger
type CfgLogger struct {
	Name        string   `yaml:"name" json:"name"`
	Level       string   `yaml:"level" json:"level"`
	OutputNames []string `yaml:"outputs" json:"outputs"`
}

// CfgOutput represents the configuration of a output
type CfgOutput map[string]string

// Name return the name of Output
func (c CfgOutput) Name() string {
	return c["name"]
}

// Type return the type of Output
func (c CfgOutput) Type() string {
	return c["type"]
}

// FormatName return the format name
func (c CfgOutput) FormatName() string {
	return c["format"]
}

// CfgFormat represents the configuration of a formatter
type CfgFormat map[string]string

// Name return the name of Format
func (c CfgFormat) Name() string {
	return c["name"]
}

// Type return the type of Format
func (c CfgFormat) Type() string {
	return c["type"]
}
