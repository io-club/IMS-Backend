package usecase

import (
	"fishnet/domain"
	"fishnet/service/device/repo"
)

var _sensorTypeUsecase domain.SensorTypeUsecase

type sensorTypeUsecase struct {
	r domain.SensorTypeRepo
}

// QuerySensorType implements domain.SensorTypeUsecase
func (*sensorTypeUsecase) QuerySensorType(sensorTypeID *int64, wordcaseID *int64, fieldName *string) ([]*domain.SensorType, error) {
	return repo.NewSensorTypeRepo().QuerySensorType(sensorTypeID, wordcaseID, fieldName)
}

// CreateSensorType implements domain.SensorTypeUsecase
func (u *sensorTypeUsecase) CreateSensorType(sensorTypes *domain.SensorType) error {
	if err := u.r.CreateSensorType([]*domain.SensorType{sensorTypes}); err != nil {
		return err
	}
	return nil
}

// DeleteSensorType implements domain.SensorTypeUsecase
func (u *sensorTypeUsecase) DeleteSensorType(sensorTypeID int64) error {
	return u.r.DeleteSensorType(sensorTypeID)
}

// MGetSensorTypes implements domain.SensorTypeUsecase
func (u *sensorTypeUsecase) MGetSensorTypes(sensorTypeIDs []int64) ([]*domain.SensorType, error) {
	return u.r.MGetSensorTypes(sensorTypeIDs)
}

// UpdateSensorType implements domain.SensorTypeUsecase
func (u *sensorTypeUsecase) UpdateSensorType(sensorTypeID int64, fieldName *string, fieldType *string) error {
	return u.r.UpdateSensorType(sensorTypeID, fieldName, fieldType)
}

func NewSensorTypeUsecase() domain.SensorTypeUsecase {
	if _sensorTypeUsecase == nil {
		_sensorTypeUsecase = &sensorTypeUsecase{
			r: repo.NewSensorTypeRepo(),
		}
	}
	return _sensorTypeUsecase
}
