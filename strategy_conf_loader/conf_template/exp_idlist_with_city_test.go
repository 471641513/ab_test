package conf_template

import (
	"ab_test/entity"
	"github.com/rs/xlog"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestNewExpIdListWithCity(t *testing.T) {
	expName := "test_exp_1"
	// 1. delete data
	err := defaultMgr.DB().Delete(&entity.ExperimentConfig{}, "name=?", expName).Error
	assert.NilError(t, err)
	now := time.Now()
	// 2. insert data for debug
	e := &entity.ExperimentConfig{
		Name:             expName,
		Status:           1,
		CityIdsStr:       `[999001,1001]`,
		ServerType:       2,
		FullOpen:         1,
		ConfStr:          `{"999001":[761]}`,
		EffectiveDateStr: now.Add(-time.Hour * 24 * 2).Format("2006-01-02 15:04:05"),
		EffectiveDays:    10,
	}
	err = defaultMgr.DB().Create(e).Error
	assert.NilError(t, err)

	exp := NewExpIdListWithCity(defaultMgr, expName)
	err = defaultMgr.Init()
	assert.NilError(t, err)
	ok, idSet := exp.GetEffectiveAndDriverSet(999001, 2)
	xlog.Info("ok=%v||idSet=%+v", ok, idSet)
	assert.Check(t, ok)
	assert.Check(t, idSet[761])

	ok, idSet = exp.GetEffectiveAndDriverSet(1001, 2)
	xlog.Info("ok=%v||idSet=%+v", ok, idSet)
	assert.Check(t, ok)
	assert.Check(t, len(idSet) == 0)
}
