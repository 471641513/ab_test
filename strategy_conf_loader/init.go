package strategy_conf_loader

import (
	"ab_test/strategy_conf_loader/conf_mgr"
	"ab_test/strategy_conf_loader/conf_template"
	"github.com/jinzhu/gorm"
)

var mgr *conf_mgr.StrategyConfigManager

var ActiveWhiteDriverScoreAdjust *conf_template.ExprIdListWithCity

var InefficientDriverScoreAdjust *conf_template.ExprIdListWithCity

func Init(db gorm.DB) (err error) {
	mgr, err = conf_mgr.NewStrategyConfigManager("", &gorm.DB{})
	if err != nil {
		return
	}

	// register exp confs
	ActiveWhiteDriverScoreAdjust = conf_template.NewExpIdListWithCity(mgr, "active_white_driver_score_adjust")

	InefficientDriverScoreAdjust = conf_template.NewExpIdListWithCity(mgr, "inefficient_driver_score_adjust")

	// init exp conf load
	err = mgr.Init()

	return
}
