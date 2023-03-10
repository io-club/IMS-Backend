package domain

import (
	"time"

	"gorm.io/gorm"
)

// int
// []int
// [][]int
// stream
type SensorData struct {
	gorm.Model
	SensorID     int64     `gorm:"index"` // 传感器 id
	SensorTypeID int64     `gorm:"index"` // 传感器字段类型 id
	Value        string    // 传感器数据
	CollectAt    time.Time `gorm:"index"` // 采集时间
}

func (SensorData) TableName() string {
	return "sensor_data"
}

type SensorDataRepo interface {
	CreateSensorData(sensorData []*SensorData) error
	DeleteSensorData(sensorDataID []int64) error
	UpdateSensorData(sensorDataID int64, SensorID, sensorTypeID *int64, sensorData *string) error
	QuerySensorData(sensorDataID *int64, sensorID *int64, sensorTypeID *int64, sensorData *string, CollectAt, beginTime, endTime *time.Time, limit, offset int) ([]*SensorData, error)
	MGetSensorDatas(sensorDataIDs []int64) ([]*SensorData, error)
}

type SensorDataUsecase interface {
	CreateSensorData(sensorData []*SensorData) error
	DeleteSensorData(sensorDataID []int64) error
	UpdateSensorData(sensorDataID int64, SensorID, sensorTypeID *int64, sensorData *string) error
	QuerySensorData(sensorDataID *int64, sensorID *int64, sensorTypeID *int64, sensorData *string, CollectAt, beginTime, endTime *time.Time, limit, offset int) ([]*SensorData, error)
	MGetSensorDatas(sensorDataIDs []int64) ([]*SensorData, error)
}
