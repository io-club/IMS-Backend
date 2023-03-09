package usecase

import (
	"fishnet/domain"
	"fishnet/service/device/repo"
)

var _deviceUsecase domain.DeviceUsecase

type deviceUsecase struct {
	r domain.DeviceRepo
}

func NewDeviceUsecase() domain.DeviceUsecase {
	if _deviceUsecase == nil {
		_deviceUsecase = &deviceUsecase{
			r: repo.NewDeviceRepo(),
		}
	}
	return _deviceUsecase
}

// CreateDevice implements domain.DeviceUsecase
func (u *deviceUsecase) CreateDevice(device *domain.Device) error {
	if err := u.r.CreateDevice([]*domain.Device{device}); err != nil {
		return err
	}
	return nil
}

// DeleteDevice implements domain.DeviceUsecase
func (u *deviceUsecase) DeleteDevice(deviceID int64) error {
	return u.r.DeleteDevice(deviceID)
}

// MGetDevices implements domain.DeviceUsecase
func (u *deviceUsecase) MGetDevices(deviceIDs []int64) ([]*domain.Device, error) {
	return u.r.MGetDevices(deviceIDs)
}

// QueryDevice implements domain.DeviceUsecase
func (u *deviceUsecase) QueryDevice(deviceID *int64, name *string, remark *string, disable *bool, limit int, offset int) ([]*domain.Device, error) {
	return u.r.QueryDevice(deviceID, name, remark, disable, limit, offset)
}

// UpdateDevice implements domain.DeviceUsecase
func (u *deviceUsecase) UpdateDevice(deviceID int64, name *string, remark *string, disable *bool) error {
	return u.r.UpdateDevice(deviceID, name, remark, disable)
}
