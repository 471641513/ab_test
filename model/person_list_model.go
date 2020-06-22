package model

import (
	"ab_test/entity"
	"github.com/jinzhu/gorm"
	"github.com/rs/xlog"
	"sync"
)

type PersonListModel struct {
	db *gorm.DB
	//	OrderPriceAdjustMap map[int32]map[int32]*OrderAdjustConf // <serverType> < city >
	lock      *sync.RWMutex
	ConfigMap map[string][]*entity.PersonListConfig
}

func NewPersonListModel(db *gorm.DB) (personListModel *PersonListModel) {

	personListModel = &PersonListModel{
		db:   db,
		lock: &sync.RWMutex{},
	}
	return
}

func (m *PersonListModel) DB() (db *gorm.DB) {
	return m.db
}

func (m *PersonListModel) UpdateConfig() (err error) {
	expConfList := make([]*entity.PersonListConfig, 0)
	exec := m.DB().Model(&entity.PersonListConfig{}).
		Where("status = 1").
		Find(&expConfList)
	if exec.Error != nil {
		return err
	}
	configMap := map[string][]*entity.PersonListConfig{}
	for _, expConf := range expConfList {
		if expConf.IsExpire() || // 过期
			expConf.IsNotStart() { // 未生效
			xlog.Info("conf=%v||experiment expire or not effective start %v, dur days = %v", expConf.ConfStr, expConf.EffectiveDateStr, expConf.EffectiveSeconds)
			continue
		}
		_, ok := configMap[expConf.Name]
		if !ok {
			configMap[expConf.Name] = make([]*entity.PersonListConfig, 0)
		}
		configMap[expConf.Name] = append(configMap[expConf.Name], expConf)
		xlog.Info("load config||name=%v||conf=%+v", expConf.Name, expConf)
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	m.ConfigMap = configMap
	return
}

func (m *PersonListModel) GetConfigListByName(name string) (tmpConfigList []*entity.PersonListConfig) {
	xlog.Debug("name=%v", name)
	m.lock.RLock()
	defer m.lock.RUnlock()
	tmpConfigList, ok := m.ConfigMap[name]
	if ok {
		configList := make([]*entity.PersonListConfig, len(tmpConfigList))
		copy(configList, tmpConfigList)
		tmpConfigList = configList
	} else {
		tmpConfigList = make([]*entity.PersonListConfig, 0)
	}
	return
}
