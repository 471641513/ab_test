package conf_mgr

import (
	"ab_test/model"
	conf_model "ab_test/person_list_loader/conf_model"
	"github.com/jinzhu/gorm"
	"github.com/rs/xlog"
	"sync"
	"time"
)

type ListMgr struct {
	lock     *sync.RWMutex
	modelMap map[string]conf_model.ListModelIFace
	*model.PersonListModel
}

func NewListMgr(db *gorm.DB) *ListMgr {
	p := &ListMgr{
		lock:            &sync.RWMutex{},
		modelMap:        map[string]conf_model.ListModelIFace{},
		PersonListModel: model.NewPersonListModel(db),
	}
	if err := p.Update(); err != nil {
		return nil
	}
	go p.LoopUpdate()
	return p
}

func (m *ListMgr) Update() (err error) {
	err = m.PersonListModel.UpdateConfig()
	if err != nil {
		return
	}
	m.lock.RLock()
	defer m.lock.RUnlock()
	for name, modelIns := range m.modelMap {
		confs := m.GetConfigListByName(name)
		if confs == nil || len(confs) == 0 {
			xlog.Debug("load config nothing name = %v", name)
		}
		modelIns.Parse(confs)
	}
	return
}

func (m *ListMgr) LoopUpdate() {
	tick := time.Tick(time.Second * time.Duration(60))
	for {
		select {
		case <-tick:
			err := m.Update()
			if err != nil {
				xlog.Error("reload error %v", err)
			} else {
				xlog.Info("reload success")
			}
		}
	}

}

func (m *ListMgr) Register(ins conf_model.ListModelIFace) (err error) {
	// 实验A
	m.lock.Lock()
	defer m.lock.Unlock()
	m.modelMap[ins.Name()] = ins
	return
}
