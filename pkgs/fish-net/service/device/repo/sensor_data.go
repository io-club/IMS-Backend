package repo

import (
	"fishnet/common/consts"
	"fishnet/domain"
	"fishnet/glb"
	"time"
)

var _sensorDataRepo domain.SensorDataRepo

type sensorDataRepo struct {
}

// CreateSensorData implements domain.SensorDataRepo
func (*sensorDataRepo) CreateSensorData(sensorData []*domain.SensorData) error {
	if err := glb.DB.Create(sensorData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSensorData implements domain.SensorDataRepo
func (*sensorDataRepo) DeleteSensorData(sensorDataID []int64) error {
	return glb.DB.Where("id in (?)", sensorDataID).Delete(&domain.SensorData{}).Error
}

// MGetSensorDatas implements domain.SensorDataRepo
func (*sensorDataRepo) MGetSensorDatas(sensorDataIDs []int64) ([]*domain.SensorData, error) {
	var res []*domain.SensorData
	if err := glb.DB.Where("id in (?)", sensorDataIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QuerySensorData implements domain.SensorDataRepo
func (*sensorDataRepo) QuerySensorData(sensorDataID *int64, sensorID *int64, sensorTypeID *int64, sensorData *string, collectAt *time.Time, beginTime *time.Time, endTime *time.Time, limit, offset int) ([]*domain.SensorData, error) {
	var total int64
	var res []*domain.SensorData
	conn := glb.DB.Model(&domain.SensorData{})
	if sensorDataID != nil && *sensorDataID > 0 {
		conn = conn.Where("id = ?", *sensorDataID)
	} else {
		if sensorID != nil && *sensorID > 0 {
			conn = conn.Where("sensor_id = ?", *sensorID)
		}
		if sensorTypeID != nil && *sensorTypeID > 0 {
			conn = conn.Where("sensor_type_id = ?", *sensorTypeID)
		}
		if sensorData != nil && *sensorData != "" {
			conn = conn.Where("value = ?", *sensorData)
		}
		if collectAt != nil && !collectAt.IsZero() {
			conn = conn.Where("collect_at = ?", *collectAt)
		}
		if beginTime != nil && !beginTime.IsZero() {
			conn = conn.Where("collect_at >= ?", *beginTime)
		}
		if endTime != nil && !endTime.IsZero() {
			conn = conn.Where("collect_at <= ?", *endTime)
		}
		if limit == 0 {
			limit = consts.DefaultLimit
		}
		conn = conn.Limit(limit).Offset(offset)
		if err := conn.Count(&total).Error; err != nil {
			return res, err
		}
	}
	if err := conn.Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateSensorData implements domain.SensorDataRepo
func (*sensorDataRepo) UpdateSensorData(sensorDataID int64, SensorID *int64, sensorTypeID *int64, sensorData *string) error {
	return glb.DB.Model(&domain.SensorData{}).Where("id = ?", sensorDataID).Updates(map[string]interface{}{
		"sensor_id":      SensorID,
		"sensor_type_id": sensorTypeID,
		"value":          sensorData,
	}).Error
}

func NewSensorDataRepo() domain.SensorDataRepo {
	if _sensorDataRepo == nil {
		_sensorDataRepo = &sensorDataRepo{}
	}
	return _sensorDataRepo
}
