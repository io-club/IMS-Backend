package usecase

import (
	"fishnet/domain"
	"fishnet/service/device/repo"
	"time"
)

var _sensorDataUsecase domain.SensorDataUsecase

type sensorDataUsecase struct {
	r domain.SensorDataRepo
}

// CreateSensorData implements domain.SensorDataUsecase
func (u *sensorDataUsecase) CreateSensorData(sensorData []*domain.SensorData) error {
	if err := u.r.CreateSensorData(sensorData); err != nil {
		return err
	}
	return nil
}

// DeleteSensorData implements domain.SensorDataUsecase
func (u *sensorDataUsecase) DeleteSensorData(sensorDataID []int64) error {
	return u.r.DeleteSensorData(sensorDataID)
}

// MGetSensorDatas implements domain.SensorDataUsecase
func (u *sensorDataUsecase) MGetSensorDatas(sensorDataIDs []int64) ([]*domain.SensorData, error) {
	return u.r.MGetSensorDatas(sensorDataIDs)
}

// QuerySensorData implements domain.SensorDataUsecase
func (u *sensorDataUsecase) QuerySensorData(sensorDataID *int64, sensorID *int64, sensorTypeID *int64, sensorData *string, collectAt *time.Time, beginTime *time.Time, endTime *time.Time, limit, offset int) ([]*domain.SensorData, error) {
	return u.r.QuerySensorData(sensorDataID, sensorID, sensorTypeID, sensorData, collectAt, beginTime, endTime, limit, offset)
}

// UpdateSensorData implements domain.SensorDataUsecase
func (u *sensorDataUsecase) UpdateSensorData(sensorDataID int64, SensorID *int64, sensorTypeID *int64, sensorData *string) error {
	return u.r.UpdateSensorData(sensorDataID, SensorID, sensorTypeID, sensorData)
}

func NewSensorDataUsecase() domain.SensorDataUsecase {
	if _sensorDataUsecase == nil {
		_sensorDataUsecase = &sensorDataUsecase{
			r: repo.NewSensorDataRepo(),
		}
	}
	return _sensorDataUsecase
}
