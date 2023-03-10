package domain

import "gorm.io/gorm"

type SensorType struct {
	gorm.Model
	WordcaseID int64  `gorm:"index"` // 词表 id
	FieldName  string // 传感器字段名称
	FieldType  string // 传感器字段类型
}

func (SensorType) TableName() string {
	return "sensor_type"
}

type SensorTypeRepo interface {
	CreateSensorType(sensorTypes []*SensorType) error
	DeleteSensorType(sensorTypeID int64) error
	UpdateSensorType(sensorTypeID int64, fieldName *string, fieldType *string) error
	QuerySensorType(sensorTypeID *int64, wordcaseID *int64, fieldName *string) ([]*SensorType, error)
	MGetSensorTypes(sensorTypeIDs []int64) ([]*SensorType, error)
}

type SensorTypeUsecase interface {
	CreateSensorType(sensorTypes *SensorType) error
	DeleteSensorType(sensorTypeID int64) error
	UpdateSensorType(sensorTypeID int64, fieldName *string, fieldType *string) error
	QuerySensorType(sensorTypeID *int64, wordcaseID *int64, fieldName *string) ([]*SensorType, error)
	MGetSensorTypes(sensorTypeIDs []int64) ([]*SensorType, error)
}
