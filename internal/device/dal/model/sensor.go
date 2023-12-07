package model

import (
	"ims-server/internal/device/job"

	"gorm.io/gorm"
)

type SensorData struct {
	gorm.Model
	Type       job.SensorType `gorm:"type:int;size:15;comment:传感器类型;not null"`
	SensorData interface{}    `gorm:"type:int;size:10;comment:传感器数据;not null"`
	TerminalID string         `gorm:"type:string;comment:传感器ID;not null"`
}

func (SensorData) TableName() string {
	return "ims_aaa_sensor"
}
