package model

import (
	"ab_test/consts"
	"ab_test/entity"
	"github.com/rs/xlog"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

type ExperimentModel struct {
	db *gorm.DB
	//	OrderPriceAdjustMap map[int32]map[int32]*OrderAdjustConf // <serverType> < city >
	lock                *sync.RWMutex
	ExperimentConfigMap map[string][]*entity.ExperimentConfig
}

func NewExperimentModel(db *gorm.DB) (
	experimentModel *ExperimentModel, err error) {

	experimentModel = &ExperimentModel{
		db:   db,
		lock: &sync.RWMutex{},
	}
	return
}

func (m *ExperimentModel) DB() (db *gorm.DB) {
	return m.db
}

func (m *ExperimentModel) Group(effectiveDate time.Time, hourGap int32, turn int32) consts.ExperimentGroup {
	t := time.Now()
	dayType := int(t.Sub(effectiveDate).Hours()) / 24 % 2
	hourType := t.Hour() / int(hourGap) % 2
	xlog.Info("dayType=%v, hourType=%v, turn=%v, hourGap=%v", dayType, hourType, turn, hourGap)
	if turn == 1 {
		if (dayType ^ hourType) == 0 {
			return consts.TREATMENT_GROUP
		}
	} else {
		if hourType == 0 {
			return consts.TREATMENT_GROUP
		}
	}
	return consts.CONTROL_GROUP
}

func (m *ExperimentModel) UpdateConfig() (err error) {
	expConfList := make([]*entity.ExperimentConfig, 0)
	exec := m.DB().Model(&entity.ExperimentConfig{}).
		Where("status = 1").
		Find(&expConfList)
	if exec.Error != nil {
		return err
	}
	experimentConfigMap := map[string][]*entity.ExperimentConfig{}
	for _, expConf := range expConfList {
		if expConf.EffectiveDate.Add(time.Duration(expConf.EffectiveDays*24)*time.Hour).Unix() < time.Now().Unix() || // 过期
			expConf.EffectiveDate.Unix() > time.Now().Unix() { // 未生效
			xlog.Info("city=%v||serverType=%v||experiment expire or not effective start %v, dur days = %v", expConf.CityIds, expConf.ServerType, expConf.EffectiveDateStr, expConf.EffectiveDays)
			continue
		}
		_, ok := experimentConfigMap[expConf.Name]
		if !ok {
			experimentConfigMap[expConf.Name] = make([]*entity.ExperimentConfig, 0)
		}
		experimentConfigMap[expConf.Name] = append(experimentConfigMap[expConf.Name], expConf)
		xlog.Info("load config||name=%v||conf=%+v", expConf.Name, expConf)
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.ExperimentConfigMap = experimentConfigMap
	return
}

func (m *ExperimentModel) GetConfigListByName(name string) (tmpConfigList []*entity.ExperimentConfig) {
	xlog.Debug("name=%v", name)
	m.lock.RLock()
	defer m.lock.RUnlock()
	tmpConfigList, ok := m.ExperimentConfigMap[name]
	if ok {
		configList := make([]*entity.ExperimentConfig, len(tmpConfigList))
		copy(configList, tmpConfigList)
		tmpConfigList = configList
	} else {
		tmpConfigList = make([]*entity.ExperimentConfig, 0)
	}
	return
}
