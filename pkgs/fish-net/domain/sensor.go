package domain

import "gorm.io/gorm"

type Sensor struct {
	gorm.Model
	DeviceID     int64 `gorm:"index"`
	SensorTypeID int `gorm:"index"`
	Stat         int   // 0: 无意义 1: 启用 2: 禁用
	Remark       string
}

func (Sensor) TableName() string {
	return "sensor"
}

type SensorRepo interface {
	CreateSensor(sensors []*Sensor) error
	DeleteSensor(sensorID int64) error
	UpdateSensor(sensorID int64, remark *string, stat *int) error
	QuerySensor(sensorID *int64, deviceID *int64, sensorTypeID *int, stat *int, limit, offset int) ([]*Sensor, error)
	MGetSensors(sensorIDs []int64) ([]*Sensor, error)
}

type SensorUsecase interface {
	CreateSensor(sensors *Sensor) error
	DeleteSensor(sensorID int64) error
	UpdateSensor(sensorID int64, remark *string, stat *int) error
	QuerySensor(sensorID *int64, deviceID *int64, sensorTypeID *int, stat *int, limit, offset int) ([]*Sensor, error)
	MGetSensors(sensorIDs []int64) ([]*Sensor, error)
	// CheckSensorExist(sensorID *int64, name *string) (bool, error)
	// FindByID(sensorID *int64) (*Sensor, error)
}
