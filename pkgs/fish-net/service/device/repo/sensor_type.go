package repo

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/glb"
)

var _sensorTypeRepo domain.SensorTypeRepo

type sensorTypeRepo struct {
}

// QuerySensorType implements domain.SensorTypeRepo
func (*sensorTypeRepo) QuerySensorType(sensorTypeID *int64, wordcaseID *int64, fieldName *string) ([]*domain.SensorType, error) {
	conn := glb.DB.Model(&domain.SensorType{})
	if sensorTypeID != nil && *sensorTypeID != 0 {
		conn = conn.Where("id = ?", *sensorTypeID)
	} else {
		if wordcaseID != nil && *wordcaseID != 0 {
			conn = conn.Where("wordcase_id = ?", *wordcaseID)
		}
		if fieldName != nil && *fieldName != "" {
			conn = conn.Where("field_name = ?", *fieldName)
		}
	}
	var res []*domain.SensorType
	if err := conn.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// CreateSensorType implements domain.SensorTypeRepo
func (*sensorTypeRepo) CreateSensorType(sensorTypes []*domain.SensorType) error {
	if err := glb.DB.Create(sensorTypes).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSensorType implements domain.SensorTypeRepo
func (*sensorTypeRepo) DeleteSensorType(sensorTypeID int64) error {
	return glb.DB.Where("id = ?", sensorTypeID).Delete(&domain.SensorType{}).Error
}

// MGetSensorTypes implements domain.SensorTypeRepo
func (*sensorTypeRepo) MGetSensorTypes(sensorTypeIDs []int64) ([]*domain.SensorType, error) {
	var res []*domain.SensorType
	if err := glb.DB.Where("id in (?)", sensorTypeIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateSensorType implements domain.SensorTypeRepo
func (*sensorTypeRepo) UpdateSensorType(sensorTypeID int64, fieldName *string, fieldType *string) error {
	conn := glb.DB.Model(&domain.SensorType{})
	params := make(map[string]interface{})
	if fieldName != nil && *fieldName != "" {
		params["field_name"] = *fieldName
	}
	if fieldType != nil && *fieldType != "" {
		params["field_type"] = *fieldType
	}
	if len(params) == 0 {
		return nil
	}
	return conn.Where("id = ?", sensorTypeID).Updates(params).Error
}

func NewSensorTypeRepo() domain.SensorTypeRepo {
	if _sensorTypeRepo == nil {
		_sensorTypeRepo = &sensorTypeRepo{}
	}
	return _sensorTypeRepo
}
