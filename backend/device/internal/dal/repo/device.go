package repo

import (
	"ims-server/device/internal/dal/model"
	ioginx "ims-server/pkg/ginx"
)

type dataRepo struct {
	ioginx.IRepo[model.Data]
}

func NewSensorRepo() *dataRepo {
	return &dataRepo{}
}
