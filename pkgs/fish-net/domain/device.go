package domain

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	GroupID  uint64 `gorm:"index"` // 设备组 id
	Name     string // 设备名称
	Remark   string // 备注
	Disabled bool   // 是否禁用
}

func (Device) TableName() string {
	return "device"
}

type DeviceRepo interface {
	CreateDevice(devices []*Device) error
	DeleteDevice(deviceID int64) error
	UpdateDevice(deviceID int64, name *string, remark *string, disable *bool) error
	QueryDevice(deviceID *int64, name *string, remark *string, disable *bool, limit, offset int) ([]*Device, error)
	MGetDevices(deviceIDs []int64) ([]*Device, error)
}

type DeviceUsecase interface {
	CreateDevice(devices *Device) error
	DeleteDevice(deviceID int64) error
	UpdateDevice(deviceID int64, name *string, remark *string, disable *bool) error
	QueryDevice(deviceID *int64, name *string, remark *string, disable *bool, limit, offset int) ([]*Device, error)
	MGetDevices(deviceIDs []int64) ([]*Device, error)
	// CheckDeviceExist(deviceID *int64, name *string) (bool, error)
	// FindByID(deviceID *int64) (*Device, error)
}
