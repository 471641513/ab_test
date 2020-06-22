package conf_decoder

import (
	"github.com/rs/xlog"
	"os"
	"testing"

	"gotest.tools/assert"
)

func TestMain(m *testing.M) {
	xlog.SetLogger(xlog.New(xlog.Config{
		Level:          0,
		Fields:         nil,
		Output:         nil,
		DisablePooling: true,
	}))
	// setup code...
	code := m.Run()
	// teardown code...
	os.Exit(code)
}

func Test_idListDecoder(t *testing.T) {
	confStr := `[761]`
	conf, err := IdListDecoder.NewConfFromStr(confStr)
	xlog.Info("confMap=%+v", conf)
	assert.NilError(t, err, "failed to decode str||err=%+v", err)
	idSet, err := IdListDecoder.GetConf(conf)
	assert.NilError(t, err, "failed to get sub conf")
	xlog.Info("idSet=%+v", idSet)
	assert.Check(t, idSet[761])
}

func Test_IdListWithCityDecoder(t *testing.T) {
	confStr := `{"999001":[761]}`
	confMap, err := IdListWithCityDecoder.NewConfFromStr(confStr)
	xlog.Info("confMap=%+v", confMap)
	assert.NilError(t, err, "failed to decode str||err=%+v", err)
	idSet, err := IdListWithCityDecoder.GetConf(confMap, "999001")
	assert.NilError(t, err, "failed to get sub conf")
	xlog.Info("idSet=%+v", idSet)
	assert.Check(t, idSet[761])
}
