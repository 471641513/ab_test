package conf_decoder

import (
	"ab_test/strategy_conf_loader/conf_mgr"
	"github.com/rs/xlog"
	"reflect"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type IdSet map[int64]bool

var IdListDecoder = &idListDecoder{}

type idListDecoder struct {
}

func (d *idListDecoder) NewConfFromStr(confStr string) (newConf conf_mgr.ExpConf, err error) {
	idSet := IdSet{}

	idList := []int64{}
	err = json.UnmarshalFromString(confStr, &idList)
	if err != nil {
		xlog.Error("failed to unmarshal idlist||err=%v||data=%s", err, confStr)
		return
	}
	for _, id := range idList {
		idSet[id] = true
	}
	newConf = idSet
	return
}

func (d idListDecoder) GetConf(conf conf_mgr.ExpConf) (idSet IdSet, err error) {
	idSet = map[int64]bool{}
	if conf == nil {
		return
	}
	var ok bool
	idSet, ok = conf.(IdSet)
	if !ok {
		//should not reach
		xlog.Fatal("subConf cant convert to idSet||type=%+v", reflect.TypeOf(conf))
		idSet = map[int64]bool{}
	}
	return
}

var IdListWithCityDecoder = &idListWithCityDecoder{
	confWithSubLevelDecoder: confWithSubLevelDecoder{
		subDecoderFunc: IdListDecoder.NewConfFromStr,
	},
}

type idListWithCityDecoder struct {
	confWithSubLevelDecoder
}

func (d idListWithCityDecoder) GetConf(confMap conf_mgr.ExpConf, subKey string) (idSet IdSet, err error) {
	idSet = map[int64]bool{}
	subConf, err := d.confWithSubLevelDecoder.getSubConf(confMap, subKey)
	if err != nil {
		xlog.Error("failed to get sub conf||err=%+v", err)
		return
	}
	idSet, err = IdListDecoder.GetConf(subConf)
	return
}
