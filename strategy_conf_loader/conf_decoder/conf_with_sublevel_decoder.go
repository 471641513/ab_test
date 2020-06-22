package conf_decoder

import (
	"ab_test/strategy_conf_loader/conf_mgr"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/rs/xlog"
	"reflect"
)

type confWithSubLevelDecoder struct {
	defaultKey     string
	subDecoderFunc func(string) (subConf conf_mgr.ExpConf, err error)
}

func (d *confWithSubLevelDecoder) NewConfFromStr(confStr string) (newConf conf_mgr.ExpConf, err error) {
	conf := map[string]*simplejson.Json{}
	err = json.UnmarshalFromString(confStr, &conf)
	if err != nil {
		xlog.Error("failed to unrmashal from string ")
		return
	}
	confMap := map[string]conf_mgr.ExpConf{}
	for k, subJson := range conf {
		subJsonBytes, err := subJson.MarshalJSON()
		if err != nil {
			xlog.Error("failed to marhsal sub json||err=%v", err)
			return confMap, err
		}
		confMap[k], err = d.subDecoderFunc(string(subJsonBytes))
	}
	newConf = confMap
	return
}

func (d *confWithSubLevelDecoder) getSubConf(conf conf_mgr.ExpConf, subKey string) (subConf conf_mgr.ExpConf, err error) {
	var confMap map[string]conf_mgr.ExpConf
	var ok bool
	if confMap, ok = conf.(map[string]conf_mgr.ExpConf); !ok {
		// should not reach
		err = fmt.Errorf("type mismatch||get=%v", reflect.TypeOf(conf))
		xlog.Fatal("cant convert to map[subkey]expConf||err=%v", err)
	} else {
		subConf = confMap[subKey]
		if subConf == nil {
			if d.defaultKey != "" {
				subConf = confMap[d.defaultKey]
			}
		}
		if subConf == nil {
			xlog.Debug("conf is nil")
			return
		}
	}
	return
}
