package conf_mgr

import (
	"ab_test/entity"
)

type ConfigMap map[int64]map[int64]*entity.ExperimentConfig // <city> <serverType>

func (m *ConfigMap) SetConfig(city, serverType int64, config *entity.ExperimentConfig) {
	_, ok := (*m)[city]
	if !ok {
		(*m)[city] = map[int64]*entity.ExperimentConfig{}
	}
	(*m)[city][serverType] = config
}

func (m *ConfigMap) GetConfig(city, serverType int64) (config *entity.ExperimentConfig) {
	_, ok := (*m)[city]
	if !ok {
		config = nil
	} else {
		config = (*m)[city][serverType]
	}
	return
}

type ExpConfigMap map[int64]map[int64]*expConfItem // <city> <serverType>

func (m *ExpConfigMap) SetConfig(city, serverType int64, config *expConfItem) {
	_, ok := (*m)[city]
	if !ok {
		(*m)[city] = map[int64]*expConfItem{}
	}
	(*m)[city][serverType] = config
}

func (m *ExpConfigMap) GetConfig(city, serverType int64) (config *expConfItem) {
	_, ok := (*m)[city]
	if !ok {
		config = nil
	} else {
		config = (*m)[city][serverType]
	}
	return
}

type ConfStrDecoder interface {
	NewConfFromStr(confStr string) (conf ExpConf, err error)
}

type Experiment interface {
	Name() string
	NewConfFromStr(confStr string) (conf ExpConf, err error)
}

type ExperimentBase struct {
	ExpName string
}

func (b *ExperimentBase) Name() string {
	return b.ExpName
}
