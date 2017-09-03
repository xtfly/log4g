package log

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"encoding/json"

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

		fmtcfg := m.config.GetCfgFormat(opcfg.FormatName)
		if fmtcfg == nil {
			return nil, lvl, fmt.Errorf("Not find format.name[%s] config", opcfg.FormatName)
		}

		fmtcreator, ok := m.formatterCreators[fmtcfg.Type]
		if !ok {
			return nil, lvl, fmt.Errorf("Not find registered format.type[%s] creator", fmtcfg.Type)
		}

		fmtt, ok := m.formats[fmtcfg.Name]
		if !ok {
			if fmtt, err = fmtcreator(fmtcfg); err != nil {
				return
			}
			m.formats[fmtcfg.Name] = fmtt
		}

		opcreator, ok := m.outputCreators[opcfg.Type]
		if !ok {
			return nil, Uninitialized, fmt.Errorf("Not find registered output.type[%s] creator", opcfg.Type)
		}

		op, ok := m.outputs[opcfg.Name]
		if !ok {
			if op, err = opcreator(opcfg); err != nil {
				return
			}
			op.SetFormatter(fmtt)
			m.outputs[opcfg.Name] = op
		}
		ops = append(ops, op)
	}
	return
}

func (m *defManager) LoadConfig(file string) error {
	m.Lock()
	defer m.Unlock()
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		err = yaml.Unmarshal(bs, m.config)
	}

	if strings.HasSuffix(file, ".json") {
		err = json.Unmarshal(bs, m.config)
	}

	if err != nil {
		return err
	}

	return fmt.Errorf("not support config file type %s", file[strings.LastIndex(file, "."):])
}

func (m *defManager) SetConfig(cfg *Config) error {
	m.Lock()
	// TODO check the output & format relationship in config
	m.config = cfg
	m.Unlock()
	return nil
}

func (m *defManager) Close() {
	m.Lock()
	for _, v := range m.outputs {
		v.Close()
	}
	m.Unlock()
}
