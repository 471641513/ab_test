package conf_mgr

import (
	"ab_test/entity"
	"github.com/rs/xlog"
	"sync"
)

type ExpConf interface {
}
type expConfItem struct {
	e    *entity.ExperimentConfig
	conf ExpConf
}

type ExpConfWrapper struct {
	lock    *sync.RWMutex
	exp     Experiment
	confMap *ExpConfigMap
}

func (w *ExpConfWrapper) updateConfList(configList []*entity.ExperimentConfig) (err error) {
	newConfMap := &ExpConfigMap{}
	for _, conf := range configList {
		xlog.Debug("cityids=%v", conf.CityIds)
		for _, cityId := range conf.CityIds {
			// 0.try decode
			expConf, err := w.exp.NewConfFromStr(conf.ConfStr)
			//err := m.model.ParseAndUpdateIdSet(cityId, conf.ServerType, conf.ConfStr)
			if err != nil {
				xlog.Error("load config error||name=%v||conf=%v||err=%v", w.exp.Name(), conf, err)
				continue
			}
			item := &expConfItem{
				e:    conf,
				conf: expConf,
			}
			newConfMap.SetConfig(cityId, conf.ServerType, item)
		}
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	w.confMap = newConfMap
	return
}

func (w *ExpConfWrapper) GetConf(city, serverType int64) (e *entity.ExperimentConfig, conf ExpConf) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	item := w.confMap.GetConfig(city, serverType)
	if item == nil {
		return
	}
	e = item.e
	conf = item.conf
	return
}
