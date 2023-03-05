package domain

import "gorm.io/gorm"

// int
// []int
// [][]int
// stream
type SensorData struct {
	gorm.Model
	SensorId uint32 `gorm:"index"`
	Value    float64
	Extra    string
}

func (SensorData) TableName() string {
	return "sensor_data"
}
