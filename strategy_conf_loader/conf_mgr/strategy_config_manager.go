package conf_mgr

import (
	"ab_test/model"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/xlog"
	"sync"
	"time"
)

type StrategyConfigManager struct {
	configModel *model.ExperimentModel
	experiments map[string]*ExpConfWrapper

	// 每个实验或者配置要在这里写清楚， 并搭配相关策略函数
	// 活跃白名单司机调权
	/*
		activeWhiteDriverScoreAdjust *ActiveWhiteDriverScoreAdjust
		// 低效黑名单司机调权
		inefficientDriverScoreAdjust *InefficientDriverScoreAdjust
	*/
}

func NewStrategyConfigManager(tableName string, db *gorm.DB) (strategyManager *StrategyConfigManager, err error) {
	// 1.0 basic experiment model
	if db == nil {
		return nil, errors.New(fmt.Sprintf("db nil"))
	}
	experimentModel, err := model.NewExperimentModel(db)
	if err != nil {
		xlog.Error("init experimentModel failed||err=%v", err)
		return
	}
	// 2.0 experiment model update
	strategyManager = &StrategyConfigManager{
		configModel: experimentModel,
		experiments: map[string]*ExpConfWrapper{},
	}
	return
}

// for debug
func (m *StrategyConfigManager) DB() *gorm.DB {
	return m.configModel.DB()
}

func (m *StrategyConfigManager) RegisterExperiment(exp Experiment) (expConfWrapper *ExpConfWrapper) {
	xlog.Info("exp registered||name=%+v", exp.Name())
	expConfWrapper = &ExpConfWrapper{
		lock:    &sync.RWMutex{},
		exp:     exp,
		confMap: &ExpConfigMap{},
	}
	m.experiments[exp.Name()] = expConfWrapper
	return
}

func (m *StrategyConfigManager) Init() (err error) {

	if len(m.experiments) == 0 {
		xlog.Warn("no exp registered")
	}

	if err = m.loadUpdate(); err != nil {
		xlog.Error("load config model error ||err=%v", err)
		return
	}
	go m.loopUpdate()
	/*
		// 3.0 strategy model update
		err = strategyManager.Update()
		if err != nil {
			xlog.Error("strategy manager load failed||err=%v", err)
			return
		}
		go strategyManager.LoopUpdate()*/
	return
}

func (m *StrategyConfigManager) loopUpdate() {
	tick := time.Tick(time.Second * time.Duration(30))
	for {
		select {
		case <-tick:
			err := m.loadUpdate()
			if err != nil {
				xlog.Error("reload strategy error ||err=%v", err)
			} else {
				xlog.Info("reload strategy success")
			}
		}
	}
}

func (m *StrategyConfigManager) loadUpdate() (err error) {
	// 0. reload exp from db
	err = m.configModel.UpdateConfig()
	if err != nil {
		xlog.Error("load config model error ||err=%v", err)
		return
	} else {
		xlog.Info("load config model success!")
	}

	// 1. load experiments
	for name, expW := range m.experiments {
		confList := m.configModel.GetConfigListByName(name)
		err = expW.updateConfList(confList)
		if err != nil {
			xlog.Error("reload strategy error ||err=%v", err)
		} else {
			xlog.Info("reload strategy success")
		}
	}
	return
}

/*

func (m *StrategyConfigManager) Update() (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 1.0 白名单司机
	m.activeWhiteDriverScoreAdjust = newActiveWhiteDriverScoreAdjust(m.lock)
	err = m.activeWhiteDriverScoreAdjust.ParseFromConfigList(m.configModel)
	if err != nil {
		xlog.Error("load config error||name=%v||err=%v", m.activeWhiteDriverScoreAdjust.name, err)
		return err
	}
	xlog.Info("load config success||name=%v", m.activeWhiteDriverScoreAdjust.name)
	// 2.0 黑名单
	m.inefficientDriverScoreAdjust = newInefficientDriverScoreAdjust(m.lock)
	err = m.inefficientDriverScoreAdjust.ParseFromConfigList(m.configModel)
	if err != nil {
		xlog.Error("load config error||name=%v||err=%v", m.inefficientDriverScoreAdjust.name, err)
		return err
	}
	xlog.Info("load config success||name=%v", m.inefficientDriverScoreAdjust.name)
	return
}

func (m *StrategyConfigManager) GetActiveWhiteScoreAdjust() *ActiveWhiteDriverScoreAdjust {
	m.lock.RLock()
	m.lock.RUnlock()
	return m.activeWhiteDriverScoreAdjust
}

func (m *StrategyConfigManager) GetInefficientWhiteScoreAdjust() *InefficientDriverScoreAdjust {
	m.lock.RLock()
	m.lock.RUnlock()
	return m.inefficientDriverScoreAdjust
}


*/
