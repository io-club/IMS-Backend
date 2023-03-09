package repo

import (
	"fishnet/common/consts"
	"fishnet/domain"
	"fishnet/glb"
)

var _deviceRepo domain.DeviceRepo

type deviceRepo struct {
}

func NewDeviceRepo() domain.DeviceRepo {
	if _deviceRepo == nil {
		_deviceRepo = &deviceRepo{}
	}
	return _deviceRepo
}

// CreateDevice implements domain.DeviceRepo
func (*deviceRepo) CreateDevice(devices []*domain.Device) error {
	if err := glb.DB.Create(devices).Error; err != nil {
		return err
	}
	return nil
}

// DeleteDevice implements domain.DeviceRepo
func (*deviceRepo) DeleteDevice(deviceID int64) error {
	return glb.DB.Where("id = ?", deviceID).Delete(&domain.Device{}).Error
}

// MGetDevices implements domain.DeviceRepo
func (*deviceRepo) MGetDevices(deviceIDs []int64) ([]*domain.Device, error) {
	var res []*domain.Device
	if err := glb.DB.Where("id in (?)", deviceIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QueryDevice implements domain.DeviceRepo
func (*deviceRepo) QueryDevice(deviceID *int64, name *string, remark *string, disable *bool, limit int, offset int) ([]*domain.Device, error) {
	var total int64
	var res []*domain.Device
	conn := glb.DB.Model(&domain.Device{})
	if deviceID != nil && *deviceID > 0 {
		conn = conn.Where("id = ?", *deviceID)
	} else {
		if name != nil && *name != "" {
			conn = conn.Where("name = ?", *name)
		}
		if disable != nil {
			conn = conn.Where("disabled = ?", *disable)
		}
		if limit == 0 {
			limit = consts.DefaultLimit
		}
		conn = conn.Limit(limit).Offset(offset)
	}
	if err := conn.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := conn.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateDevice implements domain.DeviceRepo
func (*deviceRepo) UpdateDevice(deviceID int64, name *string, remark *string, disable *bool) error {
	params := map[string]interface{}{}
	if name != nil && *name != "" {
		params["name"] = *name
	}
	if remark != nil && *remark != "" {
		params["remark"] = *remark
	}
	if disable != nil {
		params["disabled"] = *disable
	}
	return glb.DB.Model(&domain.Device{}).Where("id = ?", deviceID).Updates(params).Error
}
