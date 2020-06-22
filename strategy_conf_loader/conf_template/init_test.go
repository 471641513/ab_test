package conf_template

import (
	"ab_test/strategy_conf_loader/conf_mgr"
	"github.com/jinzhu/gorm"
	"github.com/rs/xlog"
	"os"
	"testing"
)

var defaultMgr *conf_mgr.StrategyConfigManager

func initDefaultMgr() (err error) {
	defaultMgr, err = conf_mgr.NewStrategyConfigManager("", &gorm.DB{})
	return
}

func TestMain(m *testing.M) {
	err := initDefaultMgr()
	if err != nil {
		xlog.Fatal("failed to init default Mgr")
		os.Exit(-1)
		return
	}

	// setup code...
	code := m.Run()
	// teardown code...
	os.Exit(code)
}
