package repo

import (
	"fishnet/common/consts"
	"fishnet/domain"
	"fishnet/glb"
)

var _sensorRepo domain.SensorRepo

type sensorRepo struct {
}

// CreateSensor implements domain.SensorRepo
func (*sensorRepo) CreateSensor(sensors []*domain.Sensor) error {
	if err := glb.DB.Create(sensors).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSensor implements domain.SensorRepo
func (*sensorRepo) DeleteSensor(sensorID int64) error {
	return glb.DB.Where("id = ?", sensorID).Delete(&domain.Sensor{}).Error
}

// MGetSensors implements domain.SensorRepo
func (*sensorRepo) MGetSensors(sensorIDs []int64) ([]*domain.Sensor, error) {
	var res []*domain.Sensor
	if err := glb.DB.Where("id in (?)", sensorIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QuerySensor implements domain.SensorRepo
func (*sensorRepo) QuerySensor(sensorID *int64, deviceID *int64, sensorTypeID *int, stat *int, limit int, offset int) ([]*domain.Sensor, error) {
	var total int64
	var res []*domain.Sensor
	conn := glb.DB.Model(&domain.Sensor{})
	if sensorID != nil && *sensorID > 0 {
		conn = conn.Where("id = ?", *sensorID)
	} else {
		if deviceID != nil && *deviceID > 0 {
			conn = conn.Where("device_id = ?", *deviceID)
		}
		if sensorTypeID != nil && *sensorTypeID > 0 {
			conn = conn.Where("sensor_type_id = ?", *sensorTypeID)
		}
		if stat != nil && *stat != 0 {
			conn = conn.Where("stat = ?", *stat)
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

// UpdateSensor implements domain.SensorRepo
func (*sensorRepo) UpdateSensor(sensorID int64, remark *string, stat *int) error {
	var sensor domain.Sensor
	if err := glb.DB.Where("id = ?", sensorID).First(&sensor).Error; err != nil {
		return err
	}
	if remark != nil {
		sensor.Remark = *remark
	}
	if stat != nil && *stat != 0 {
		sensor.Stat = *stat
	}
	if err := glb.DB.Save(&sensor).Error; err != nil {
		return err
	}
	return nil
}

func NewSensorRepo() domain.SensorRepo {
	if _sensorRepo == nil {
		_sensorRepo = &sensorRepo{}
	}
	return _sensorRepo
}
