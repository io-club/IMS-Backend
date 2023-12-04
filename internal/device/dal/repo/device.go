package repo

import (
	"ims-server/internal/device/dal/model"
	ioginx "ims-server/pkg/ginx"
)



type sensorRepo struct {
	ioginx.IRepo[model.Sensor]
}

func NewSensorRepo() *sensorRepo {
	return &sensorRepo{}
}
