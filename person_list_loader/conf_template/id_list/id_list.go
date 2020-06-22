package id_list

import (
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
	"reflect"
)

type IdListConfig struct {
	IdMap map[int64]bool
}

func (m *IdListConfig) ParseConfig(confStr string) (confMap map[int64]map[int64]bool, err error) {
	return map[int64]map[int64]bool{0: {0: true}}, nil
}

func (m *IdListConfig) ParseList(listStr string) (err error) {
	list := []int64{}
	err = json.UnmarshalFromString(listStr, &list)
	if err != nil {
		return
	}
	m.IdMap = map[int64]bool{}
	for _, id := range list {
		m.IdMap[id] = true
	}
	return
}

func (m *IdListConfig) Merge(conf interface{}) (err error) {
	if reflect.TypeOf(conf) != reflect.TypeOf(&IdListConfig{}) {
		return errors.New(fmt.Sprintf("type error!"))
	}
	var c *IdListConfig = conf.(*IdListConfig)
	for id := range c.IdMap {
		m.IdMap[id] = true
	}
	return nil
}

func (m *IdListConfig) IdInList(did int64) bool {
	if r, ok := m.IdMap[did]; ok && r {
		return true
	}
	return false
}
