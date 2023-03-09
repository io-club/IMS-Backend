package http

import (
	"fishnet/domain"
	"fishnet/service/device/usecase"
)

var _deviceUsecase domain.DeviceUsecase

func init() {
	_deviceUsecase = usecase.NewDeviceUsecase()
}
