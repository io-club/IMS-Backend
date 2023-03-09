package domain

import "gorm.io/gorm"

type Sensor struct {
	gorm.Model
	DeviceID     uint32 `gorm:"index"`
	SensorTypeID uint32 `gorm:"index"`
	Disabled     bool
	Remark       string
}

func (Sensor) TableName() string {
	return "sensor"
}
