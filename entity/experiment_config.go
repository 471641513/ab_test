package entity

import (
	"ab_test/consts"
	"github.com/rs/xlog"
	"time"

	json "github.com/json-iterator/go"
)

type ExperimentConfig struct {
	Id               int64     `gorm:"column:id" json:"id"`
	Name             string    `gorm:"column:name" json:"name"`
	Status           int32     `gorm:"column:status" json:"status"`
	CityIdsStr       string    `gorm:"column:city_ids_str" json:"city_ids_str"`
	ServerType       int64     `gorm:"column:server_type" json:"server_type"`
	CityIds          []int64   `gorm:"-" json:"city_ids"`
	FullOpen         int32     `gorm:"column:full_open" json:"full_open"`
	Turn             int32     `gorm:"column:turn" json:"turn"`
	HourGap          int32     `gorm:"column:hour_gap" json:"hour_gap"`
	ConfStr          string    `gorm:"column:conf_str" json:"conf_str"`
	EffectiveDateStr string    `gorm:"column:effective_date_str" json:"effective_date_str"`
	EffectiveDate    time.Time `gorm:"-" json:"effective_date"`
	EffectiveDays    int64     `gorm:"column:effective_days" json:"effective_days"`
}

func (m *ExperimentConfig) AfterFind() (err error) {
	defer func() {
		if err != nil {
			xlog.Error("ParseErr %v", err)
		}
	}()
	m.CityIds = []int64{}
	if m.CityIdsStr != "" {
		err = json.UnmarshalFromString(m.CityIdsStr, &m.CityIds)
		if err != nil {
			return err
		}
	}
	effTime, err := time.ParseInLocation("2006-01-02 15:04:05", m.EffectiveDateStr, time.Local)
	if err != nil {
		return err
	}
	m.EffectiveDate = effTime
	return
}

func (m *ExperimentConfig) Group() consts.ExperimentGroup {
	t := time.Now()
	effectiveDate := m.EffectiveDate
	hourGap := m.HourGap
	turn := m.Turn
	dayType := int(t.Sub(effectiveDate).Hours()) / 24 % 2
	if hourGap == 0 {
		return consts.CONTROL_GROUP
	}
	hourType := t.Hour() / int(hourGap) % 2
	xlog.Info("dayType=%v, hourType=%v, turn=%v, hourGap=%v", dayType, hourType, turn, hourGap)
	if turn == 1 {
		if (dayType ^ hourType) == 0 {
			return consts.TREATMENT_GROUP
		}
	} else {
		if hourType == 0 {
			return consts.TREATMENT_GROUP
		}
	}
	return consts.CONTROL_GROUP
}

func (m *ExperimentConfig) IsFullOpen() bool {
	return m.FullOpen == 1
}

func (m *ExperimentConfig) Effective() bool {
	if m.Group().IsTreatmentGroup() || m.IsFullOpen() {
		return true
	} else {
		return false
	}
}
