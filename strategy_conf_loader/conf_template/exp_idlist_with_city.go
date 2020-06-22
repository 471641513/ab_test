package conf_template

import (
	"ab_test/strategy_conf_loader/conf_decoder"
	"ab_test/strategy_conf_loader/conf_mgr"
	"fmt"
)

type ExprIdListWithCity struct {
	conf_mgr.ExperimentBase
	conf_mgr.ConfStrDecoder
	wrapper *conf_mgr.ExpConfWrapper
}

func NewExpIdListWithCity(mgr *conf_mgr.StrategyConfigManager, name string) (exp *ExprIdListWithCity) {
	exp = &ExprIdListWithCity{
		ExperimentBase: conf_mgr.ExperimentBase{
			ExpName: name,
		},
		ConfStrDecoder: conf_decoder.IdListWithCityDecoder,
	}
	exp.wrapper = mgr.RegisterExperiment(exp)
	return
}

func (m *ExprIdListWithCity) subKey(city int64) string {
	return fmt.Sprintf("%d", city)
}

func (m *ExprIdListWithCity) GetEffectiveAndDriverSet(city, serverType int64) (ok bool, idSet conf_decoder.IdSet) {
	config, confMap := m.wrapper.GetConf(city, serverType)
	if config == nil {
		ok = false
	} else {
		ok = config.Effective()
	}
	idSet, _ = conf_decoder.IdListWithCityDecoder.GetConf(confMap, m.subKey(city))
	return
}
