package http

import (
	"fishnet/domain"
	"fishnet/service/device/usecase"
)

var _deviceUsecase domain.DeviceUsecase
var _sensorUsecase domain.SensorUsecase

func init() {
	_deviceUsecase = usecase.NewDeviceUsecase()
	_sensorUsecase = usecase.NewSensorUsecase()
}
