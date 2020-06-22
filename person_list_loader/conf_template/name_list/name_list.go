package name_list

import json "github.com/json-iterator/go"

type NameListConfig struct {
	NameMap  map[string]bool
	NameList []string
}

func (m *NameListConfig) ParseConfig(confStr string) (confMap map[int64]map[int64]bool, err error) {
	return map[int64]map[int64]bool{0: {0: true}}, nil
}

func (m *NameListConfig) ParseList(listStr string) (err error) {
	list := []string{}
	err = json.UnmarshalFromString(listStr, &list)
	if err != nil {
		return
	}
	m.NameMap = map[string]bool{}
	m.NameList = list
	for _, id := range m.NameList {
		m.NameMap[id] = true
	}
	return
}

func (m *NameListConfig) NameInList(name string) bool {
	if r, ok := m.NameMap[name]; ok && r {
		return true
	}
	return false
}
