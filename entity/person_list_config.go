package entity

import (
	"github.com/rs/xlog"
	"time"
)

type PersonListConfig struct {
	Id               int64     `gorm:"column:id" json:"id"`
	Name             string    `gorm:"column:name" json:"name"`
	Status           int32     `gorm:"column:status" json:"status"`
	FullOpen         int32     `gorm:"column:full_open" json:"full_open"`
	ConfStr          string    `gorm:"column:conf_str" json:"conf_str"`
	ListStr          string    `gorm:"column:list_str" json:"list_str"`
	EffectiveDateStr string    `gorm:"column:effective_date_str" json:"effective_date_str"`
	EffectiveDate    time.Time `gorm:"-" json:"effective_date"`
	EffectiveSeconds int64     `gorm:"column:effective_seconds" json:"effective_seconds"`
}

func (m *PersonListConfig) AfterFind() (err error) {
	defer func() {
		if err != nil {
			xlog.Error("ParseErr %v", err)
		}
	}()
	effTime, err := time.ParseInLocation("2006-01-02 15:04:05", m.EffectiveDateStr, time.Local)
	if err != nil {
		return err
	}
	m.EffectiveDate = effTime
	return
}

func (m *PersonListConfig) IsFullOpen() bool {
	return m.FullOpen == 1
}

func (m *PersonListConfig) IsExpire() bool {
	if m.EffectiveDate.Add(time.Duration(m.EffectiveSeconds)*time.Second).Unix() < time.Now().Unix() {
		return true
	}
	return false
}

func (m *PersonListConfig) IsNotStart() bool {
	if m.EffectiveDate.Unix() > time.Now().Unix() { // 未生效
		return true
	}
	return false
}
