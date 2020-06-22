package conf_model

import (
	"ab_test/entity"
	"ab_test/person_list_loader/conf_template"
	"sync"
)

type ListModelIFace interface {
	Clear()
	Update(cityId, serverType int64, template conf_template.ListTemplate) (err error)
	Get(cityId, serverType int64) conf_template.ListTemplate
	Parse([]*entity.PersonListConfig)
	Name() string
}

type ListModel struct {
	lock    *sync.RWMutex
	confMap map[int64]map[int64]conf_template.ListTemplate
	name    string
}

func (m *ListModel) Name() string {
	return m.name
}

func NewListModel(name string) *ListModel {
	return &ListModel{
		lock:    &sync.RWMutex{},
		confMap: map[int64]map[int64]conf_template.ListTemplate{},
		name:    name,
	}
}

func (m *ListModel) Clear() {
	m.lock.Lock()
	m.confMap = map[int64]map[int64]conf_template.ListTemplate{}
	m.lock.Unlock()
}

func (m *ListModel) Update(cityId, serverType int64, template conf_template.ListTemplate) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.confMap[cityId]; !ok {
		m.confMap[cityId] = map[int64]conf_template.ListTemplate{}
	}
	if _, ok := m.confMap[cityId][serverType]; !ok {
		m.confMap[cityId][serverType] = template
		return
	}
	err = m.confMap[cityId][serverType].Merge(template)
	return
}

func (m *ListModel) Get(cityId, serverType int64) conf_template.ListTemplate {
	m.lock.RLock()
	serverMap, ok := m.confMap[cityId]
	if !ok {
		return nil
	}
	conf, ok := serverMap[serverType]
	if !ok || conf == nil {
		return nil
	}
	m.lock.RUnlock()
	return conf
}
