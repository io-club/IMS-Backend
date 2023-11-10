package usecase

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/device/repo"
)

var _sensorUsecase domain.SensorUsecase

type sensorUsecase struct {
	r domain.SensorRepo
}

// CreateSensor implements domain.SensorUsecase
func (u *sensorUsecase) CreateSensor(sensors *domain.Sensor) error {
	if err := u.r.CreateSensor([]*domain.Sensor{sensors}); err != nil {
		return err
	}
	return nil
}

// DeleteSensor implements domain.SensorUsecase
func (u *sensorUsecase) DeleteSensor(sensorID int64) error {
	return repo.NewSensorRepo().DeleteSensor(sensorID)
}

// MGetSensors implements domain.SensorUsecase
func (u *sensorUsecase) MGetSensors(sensorIDs []int64) ([]*domain.Sensor, error) {
	return u.r.MGetSensors(sensorIDs)
}

// QuerySensor implements domain.SensorUsecase
func (u *sensorUsecase) QuerySensor(sensorID *int64, deviceID *int64, wordcaseID *int, stat *int, limit int, offset int) ([]*domain.Sensor, error) {
	return u.r.QuerySensor(sensorID, deviceID, wordcaseID, stat, limit, offset)
}

// UpdateSensor implements domain.SensorUsecase
func (u *sensorUsecase) UpdateSensor(sensorID int64, remark *string, stat *int) error {
	return u.r.UpdateSensor(sensorID, remark, stat)
}

func NewSensorUsecase() domain.SensorUsecase {
	if _sensorUsecase == nil {
		_sensorUsecase = &sensorUsecase{
			r: repo.NewSensorRepo(),
		}
	}
	return _sensorUsecase
}
