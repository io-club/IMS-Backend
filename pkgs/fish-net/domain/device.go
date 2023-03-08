package domain

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	GroupID      uint64 `gorm:"index"` // 设备组 id
	Name         uint64 // 设备名称
	SendorTypeID uint8  // 设备类型
	Remark       string // 备注
	Disabled     bool   // 是否禁用
}

func (Device) TableName() string {
	return "device"
}

type DeviceRepo interface {
	CreateDevice(devices []*Device) error
	DeleteDevice(deviceID int64) error
	UpdateDevice(deviceID int64, name *string, remark *string) error
	QueryDevice(deviceID *int64, name *string, remark *string, limit, offset int) ([]*Device, int64, error)
	MGetDevices(deviceIDs []int64) ([]*Device, error)
}

type DeviceUsecase interface {
	CreateDevice(devices []*Device) error
	DeleteDevice(deviceID int64) error
	UpdateDevice(deviceID int64, name *string, remark *string) error
	QueryDevice(deviceID *int64, name *string, remark *string, limit, offset int) ([]*Device, int64, error)
	MGetDevices(deviceIDs []int64) ([]*Device, error)
	CheckDeviceExist(deviceID *int64, name *string) (bool, error)
	FindByID(deviceID *int64) (*Device, error)
}
