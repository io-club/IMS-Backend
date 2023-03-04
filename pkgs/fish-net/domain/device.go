package domain

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Id      uint64 // id
	GroupId uint64 `gorm:"index"` // 设备组 id
	Name    uint64 // 设备名？
	Type    uint8  // 设备类型
}

type Sensor struct {
	gorm.Model
	// 设备 id
	// 传感器与设备是一一对应的吗？
	DeviceId uint32 `gorm:"index"`
	Type     uint8
}

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
