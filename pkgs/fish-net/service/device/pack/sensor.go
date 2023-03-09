package pack

import (
	"fishnet/domain"
	"fishnet/service/common"
)

type CreateSensorRequest struct {
	SensorTypeID int    `json:"sensorTypeId"` // 传感器类型
	Remark       string `json:"remark"`       // 备注
	Stat         int    `json:"stat"`         // 0: 无意义 1: 启用 2: 禁用
}

type UpdateSensorRequest struct {
	Remark string `json:"remark"` // 备注
	Stat   int    `json:"stat"`   // 0: 无意义 1: 启用 2: 禁用
}

type QuerySensorRequest struct {
	SensorType int `json:"sensorType"` // 传感器类型
	Stat       int `json:"stat"`       // 0: 无意义 1: 启用 2: 禁用
	common.PageRequest
}

type QuerySensorResponseValue struct {
	ID         int64  `json:"id"`         // 传感器 id
	SensorType int    `json:"sensorType"` // 传感器类型
	Remark     string `json:"remark"`     // 备注
	Stat       int    `json:"stat"`       // 0: 无意义 1: 启用 2: 禁用
}

type QuerySensorResponseEntries map[int64]QuerySensorResponseValue

func Value(s *domain.Sensor) QuerySensorResponseValue {
	return QuerySensorResponseValue{
		ID:         int64(s.ID),
		SensorType: int(s.SensorTypeID),
		Remark:     s.Remark,
		Stat:       s.Stat,
	}
}

func Values(ss []*domain.Sensor) []QuerySensorResponseValue {
	res := make([]QuerySensorResponseValue, len(ss))
	for i, s := range ss {
		res[i] = Value(s)
	}
	return res
}

func Entrys(ss []*domain.Sensor) QuerySensorResponseEntries {
	res := make(QuerySensorResponseEntries, len(ss))
	for _, s := range ss {
		res[int64(s.DeviceID)] = Value(s)
	}
	return res
}
