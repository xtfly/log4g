package log

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"encoding/json"

	"path"

	"gopkg.in/yaml.v2"
)

type defManager struct {
	sync.Mutex
	formatterCreators map[string]FormatterFuncCreator // key: type
	outputCreators    map[string]OutputFuncCreator    // key: type
	formats           map[string]Formatter            // key: name
	outputs           map[string]Output               // key: name
	config            *Config
}

func newManager() Manager {
	return &defManager{
		formatterCreators: make(map[string]FormatterFuncCreator),
		outputCreators:    make(map[string]OutputFuncCreator),
		formats:           make(map[string]Formatter),
		outputs:           make(map[string]Output),
		config:            &Config{},
	}
}

func (m *defManager) RegisterFormatterCreator(stype string, f FormatterFuncCreator) {
	m.Lock()
	m.formatterCreators[stype] = f
	m.Unlock()
}

func (m *defManager) RegisterOutputCreator(stype string, o OutputFuncCreator) {
	m.Lock()
	m.outputCreators[stype] = o
	m.Unlock()
}

func (m *defManager) GetLoggerOutputs(name string) (ops []Output, lvl Level, err error) {
	lvl = Uninitialized
	m.Lock()
	defer m.Unlock()
	lc := m.config.GetCfgLogger(name)
	if lc == nil {
		return nil, lvl, fmt.Errorf("Not find logger.name[%s] config", name)
	}
	lvl = LevelFrom(lc.Level)

	for _, opid := range lc.OutputNames {
		opcfg := m.config.GetCfgOutput(opid)
		if opcfg == nil {
			return nil, lvl, fmt.Errorf("Not find output.name[%s] config", opid)
		}

		fmtcfg := m.config.GetCfgFormat(opcfg.FormatName())
		if fmtcfg == nil {
			return nil, lvl, fmt.Errorf("Not find format.name[%s] config", opcfg.FormatName())
		}

		fmtcreator, ok := m.formatterCreators[fmtcfg.Type()]
		if !ok {
			return nil, lvl, fmt.Errorf("Not find registered format.type[%s] creator", fmtcfg.Type())
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
			return nil, Uninitialized, fmt.Errorf("Not find registered output.type[%s] creator", opcfg.Type())
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
	m.Lock()
	defer m.Unlock()

	cfg := &Config{}
	if ext == "yaml" || ext == "yml" {
		err = yaml.Unmarshal(bs, cfg)
	} else if ext == "json" {
		err = json.Unmarshal(bs, cfg)
	}
	return m.setConfig(cfg)
}

func (m *defManager) setConfig(cfg *Config) error {
	// TODO check the output & format relationship in config
	m.config = cfg
	return nil
}

func (m *defManager) SetConfig(cfg *Config) error {
	m.Lock()
	defer m.Unlock()
	return m.setConfig(cfg)
}

func (m *defManager) Close() {
	m.Lock()
	for _, v := range m.outputs {
		v.Close()
	}
	m.Unlock()
}
