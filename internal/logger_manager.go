package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"sync"

	"encoding/json"

	"path"

	"github.com/xtfly/log4g/api"
)

type configNotification interface {
	notify()
}

type defManager struct {
	sync.RWMutex
	formatterCreators map[string]api.FormatterFuncCreator // key: type
	outputCreators    map[string]api.OutputFuncCreator    // key: type
	formats           map[string]api.Formatter            // key: name
	outputs           map[string]api.Output               // key: name
	config            *api.Config
	cfgNotifications  []configNotification
}

func newManager() api.Manager {
	return &defManager{
		formatterCreators: make(map[string]api.FormatterFuncCreator),
		outputCreators:    make(map[string]api.OutputFuncCreator),
		formats:           make(map[string]api.Formatter),
		outputs:           make(map[string]api.Output),
		config:            &api.Config{},
	}
}

func (m *defManager) RegisterFormatterCreator(stype string, f api.FormatterFuncCreator) {
	m.Lock()
	m.formatterCreators[stype] = f
	m.Unlock()
}

func (m *defManager) RegisterOutputCreator(stype string, o api.OutputFuncCreator) {
	m.Lock()
	m.outputCreators[stype] = o
	m.Unlock()
}

func (m *defManager) GetLoggerOutputs(name string) (ops []api.Output, lvl api.Level, err error) {
	lvl = api.Uninitialized
	m.Lock()
	defer m.Unlock()
	lc := m.config.GetCfgLogger(name)
	if lc == nil {
		return nil, lvl, fmt.Errorf("not find logger.name[%s] config", name)
	}
	lvl = api.LevelFrom(lc.Level)

	for _, opid := range lc.OutputNames {
		opcfg := m.config.GetCfgOutput(opid)
		if opcfg == nil {
			return nil, lvl, fmt.Errorf("not find output.name[%s] config", opid)
		}

		fmtcfg := m.config.GetCfgFormat(opcfg.FormatName())
		if fmtcfg == nil {
			return nil, lvl, fmt.Errorf("not find format.name[%s] config", opcfg.FormatName())
		}

		fmtcreator, ok := m.formatterCreators[fmtcfg.Type()]
		if !ok {
			return nil, lvl, fmt.Errorf("not find registered format.type[%s] creator", fmtcfg.Type())
		}

		fmtt, ok := m.formats[fmtcfg.Name()]
		if !ok {
			if fmtt, err = fmtcreator(fmtcfg); err != nil {
				return
			}
			m.formats[fmtcfg.Name()] = fmtt
		}

		opcreator, ok := m.outputCreators[opcfg.Type()]
		if !ok {
			return nil, api.Uninitialized, fmt.Errorf("not find registered output.type[%s] creator", opcfg.Type())
		}

		op, ok := m.outputs[opcfg.Name()]
		if !ok {
			if op, err = opcreator(opcfg); err != nil {
				return
			}
			op.SetFormatter(fmtt)
			m.outputs[opcfg.Name()] = op
		}
		ops = append(ops, op)
	}
	return
}

func (m *defManager) LoadConfigFile(file string) error {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	ext := path.Ext(file)
	return m.LoadConfig(bs, ext[1:])
}

func (m *defManager) LoadConfig(bs []byte, ext string) (err error) {
	ext = strings.ToLower(ext)
	err = fmt.Errorf("not support config file type %s", ext)

	cfg := &api.Config{}
	if ext == "yaml" || ext == "yml" {
		err = yaml.Unmarshal(bs, cfg)
	} else if ext == "json" {
		err = json.Unmarshal(bs, cfg)
	}

	return m.setConfig(cfg)
}

func (m *defManager) setConfig(cfg *api.Config) error {
	if err := m.validateConfig(cfg); err != nil {
		return err
	}

	m.Lock()
	m.config = cfg
	m.Unlock()
	for _, cn := range m.cfgNotifications {
		cn.notify()
	}
	return nil
}

func (m *defManager) validateConfig(cfg *api.Config) (err error) {
	// check the output & format relationship in config

	formats := make(map[string]api.CfgFormat)
	for _, f := range cfg.Formats {
		if _, ok := formats[f.Name()]; ok {
			err = fmt.Errorf("duplication formatter[%s] config", f.Name())
			return
		}
		formats[f.Name()] = f
	}

	outputs := make(map[string]api.CfgOutput)
	for _, o := range cfg.Outputs {
		if _, ok := outputs[o.Name()]; ok {
			err = fmt.Errorf("duplication output[%s] config", o.Name())
			return
		}
		if _, ok := formats[o.FormatName()]; !ok {
			err = fmt.Errorf("not found format[%s] for output[%s] ", o.FormatName(), o.Name())
			return
		}
		outputs[o.Name()] = o
	}

	loggers := make(map[string]api.CfgLogger)
	for _, l := range cfg.Loggers {
		if _, ok := loggers[l.Name]; ok {
			err = fmt.Errorf("duplication logger[%s] config", l.Name)
			return
		}

		for _, outputName := range l.OutputNames {
			if _, ok := outputs[outputName]; !ok {
				err = fmt.Errorf("not found output[%s] for logger[%s] ", outputName, l.Name)
				return
			}
		}
	}

	return
}

func (m *defManager) SetConfig(cfg *api.Config) error {
	return m.setConfig(cfg)
}

func (m *defManager) Close() {
	m.Lock()
	for _, v := range m.outputs {
		v.Close()
	}
	m.Unlock()
}

func (m *defManager) addConfigNotify(cn configNotification) {
	m.cfgNotifications = append(m.cfgNotifications, cn)
}
