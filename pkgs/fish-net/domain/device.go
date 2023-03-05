package domain

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Id      uint64 // id
	GroupId uint64 `gorm:"index"` // 设备组 id
	Name    uint64 // 设备名？
	Type    uint8  // 设备类型
}

func (Device) TableName() string {
	return "device"
}
