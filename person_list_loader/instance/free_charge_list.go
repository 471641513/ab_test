package instance

import (
	"ab_test/entity"
	"ab_test/person_list_loader/conf_model"
	"ab_test/person_list_loader/conf_template/id_list"
)

type FreeCharge struct {
	*conf_model.ListModel
}

func NewFreeCharge(name string) conf_model.ListModelIFace {
	return &FreeCharge{conf_model.NewListModel(name)}
}

func (m *FreeCharge) Parse(confList []*entity.PersonListConfig) {
	for _, conf := range confList {
		ins := &id_list.IdListConfig{}
		cityServerMap, err := ins.ParseConfig(conf.ConfStr)
		if err != nil {
			return
		}
		err = ins.ParseList(conf.ListStr)
		if err != nil {
			return
		}
		for cityId, serverMap := range cityServerMap {
			for serverType := range serverMap {
				err = m.Update(cityId, serverType, ins)
				if err != nil {
					return
				}
			}
		}
	}
	return
}
