package pack

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/common"
)

type CreateDeviceRequest struct {
	Name   string // 设备名称
	Remark string `json:"remark"` // 备注
}

type UpdateDeviceRequest struct {
	Name    string `json:"name"`    // 设备名称
	Remark  string `json:"remark"`  // 备注
	Disable bool   `json:"disable"` // 是否禁用
}

type QueryDeviceRequest struct {
	Name    string `json:"name"`    // 设备名称
	Remark  string `json:"remark"`  // 备注
	Disable bool   `json:"disable"` // 是否禁用
	common.PageRequest
}

type QueryDeviceResponse struct {
	ID      int64  `json:"id"`      // 设备 id
	Name    string `json:"name"`    // 设备名称
	Remark  string `json:"remark"`  // 备注
	Disable bool   `json:"disable"` // 是否禁用
}

func Device(d *domain.Device) QueryDeviceResponse {
	return QueryDeviceResponse{
		ID:      int64(d.ID),
		Name:    d.Name,
		Remark:  d.Remark,
		Disable: d.Disabled,
	}
}

func Devices(ds []*domain.Device) []QueryDeviceResponse {
	res := make([]QueryDeviceResponse, len(ds))
	for i, d := range ds {
		res[i] = Device(d)
	}
	return res
}
