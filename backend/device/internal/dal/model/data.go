package model

import (
	"gorm.io/gorm"
	"ims-server/device/internal/job"
)

type Data struct {
	gorm.Model
	Type       job.SensorType `gorm:"type:string;size:15;comment:传感器类型;not null"`
	TerminalID string         `gorm:"type:string;comment:终端 ID;not null"`
	SensorData interface{}    `gorm:"type:string;comment:传感器数据 not null"`
}

func (Data) TableName() string {
	return "ims_mqtt_SensorDate"
}
