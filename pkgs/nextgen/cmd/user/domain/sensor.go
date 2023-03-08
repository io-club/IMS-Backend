package domain

import "gorm.io/gorm"

type Sensor struct {
	gorm.Model
	// 设备 id
	// 传感器与设备是一一对应的吗？
	DeviceId uint32 `gorm:"index"`
	Type     uint8
}

func (Sensor) TableName() string {
	return "sensor"
}
