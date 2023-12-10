package model

import (
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	SensorType string `gorm:"type:char(15);comment:传感器类型;not null"`
	TerminalID string `gorm:"type:string;comment:终端 ID;not null"`
	Data       string `gorm:"type:string;comment:传感器数据 not null"`
}

func (Data) TableName() string {
	return "ims_mqtt_date"
}
