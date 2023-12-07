package model

import (
	"gorm.io/gorm"
	"ims-server/internal/device/job"
)

type SensorData struct {
	gorm.Model
	Type       job.SensorType `gorm:"type:cstring;size:15;comment:传感器类型;not null"`
	TerminalID string         `gorm:"type:string;comment:终端ID;not null"`
	SensorData interface{}    `gorm:"type:string;comment:传感器数据not null"`
}

func (SensorData) TableName() string {
	return "ims_mqtt_SensorDate"
}
