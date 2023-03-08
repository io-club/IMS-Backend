package domain

import "gorm.io/gorm"

type SensorType struct {
	gorm.Model
	FieldName string
	FieldType string
}

func (SensorType) TableName() string {
	return "sensor_type"
}
