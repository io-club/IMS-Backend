package repo

import (
	"context"
	"github.com/bytedance/sonic"
	"ims-server/microservices/device/internal/dal/model"
	"ims-server/microservices/device/internal/job"
	ioginx "ims-server/pkg/ginx"
)

type dataRepo struct {
	ioginx.IRepo[model.Data]
}

func NewDataRepo() *dataRepo {
	return &dataRepo{}
}

func (d *dataRepo) CreateFromRawData(ctx context.Context, sensorType job.SensorType, rawData *job.UartMsg) error {
	sensorData, err := sonic.Marshal(rawData.Sensor)
	if err != nil {
		return err
	}
	data := &model.Data{
		SensorType: sensorType.ToString(),
		TerminalID: rawData.TerminalID,
		Data:       string(sensorData),
	}
	err = d.Create(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
