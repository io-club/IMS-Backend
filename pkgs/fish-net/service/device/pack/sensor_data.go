package pack

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/common"
	"time"
)

type CreateSensorDataRequest struct {
	SensorID     *int64    `json:"sensorId"`     // 传感器 id
	SensorTypeID *int64    `json:"sensorTypeId"` // 传感器字段类型 id
	Value        string    `json:"value"`        // 传感器数据内容
	CollectTime  time.Time `json:"collectTime"`  // 采集时间
}

type UpdateSensorDataRequest struct {
	Value *string `json:"value"` // 传感器数据内容
}

type QuerySensorDataRequest struct {
	SensorId   *int64     `json:"sensorId"`   // 传感器 id
	WordcaseID *int64     `json:"wordcaseId"` // 传感器类型
	Value      *string    `json:"value"`      // 传感器数据内容
	BeginTime  *time.Time `json:"beginTime"`  // 开始时间
	EndTime    *time.Time `json:"endTime"`    // 结束时间
	common.PageRequest
}

type QuerySensorDataResponse struct {
	ID           int64  `json:"id"`           // 传感器数据 id
	SensorID     int64  `json:"sensorId"`     // 传感器 id
	SensorTypeID int64  `json:"sensorTypeId"` // 传感器类型 ID
	FieldName    string `json:"fieldName"`    // 传感器类型字段名
	Value        string `json:"value"`        // 传感器数据内容
	CollectTime  string `json:"collectTime"`  // 采集时间
}

func SensorData(s *domain.SensorData, sensorTypeID2SensorType map[int64]*domain.SensorType) QuerySensorDataResponse {
	return QuerySensorDataResponse{
		ID:           int64(s.ID),
		SensorID:     s.SensorID,
		SensorTypeID: s.SensorTypeID,
		FieldName:    sensorTypeID2SensorType[s.SensorTypeID].FieldName,
		Value:        s.Value,
		CollectTime:  s.CollectAt.Format("2006-01-02 15:04:05"),
	}
}

func SensorDatas(ss []*domain.SensorData, sensorTypeID2SensorType map[int64]*domain.SensorType) []QuerySensorDataResponse {
	res := make([]QuerySensorDataResponse, len(ss))
	for i, s := range ss {
		res[i] = SensorData(s, sensorTypeID2SensorType)
	}
	return res
}
