package pack

import "fishnet/domain"

type CreateSensorTypeRequest struct {
	WordcaseID int64  `json:"wordcaseId"` // 词表 id
	FieldType  string `json:"fieldType"`  // 传感器字段类型
	FieldName  string `json:"fieldName"`  // 传感器字段名称
}

type QuerySensorTypeRequest struct {
	WordcaseID int64  `json:"wordcaseId"` // 词表 id
	FieldType  string `json:"fieldType"`  // 传感器字段类型
	FieldName  string `json:"fieldName"`  // 传感器字段名称
}

type QuerySensorTypeResponse struct {
	ID        int64  `json:"id"`        // 传感器类型 id
	FieldType string `json:"fieldType"` // 传感器字段类型
	FieldName string `json:"fieldName"` // 传感器字段名称
}

func SensorType(s *domain.SensorType) QuerySensorTypeResponse {
	return QuerySensorTypeResponse{
		ID:        int64(s.ID),
		FieldType: s.FieldType,
		FieldName: s.FieldName,
	}
}

func SensorTypes(ss []*domain.SensorType) []QuerySensorTypeResponse {
	var res []QuerySensorTypeResponse
	for _, s := range ss {
		res = append(res, SensorType(s))
	}
	return res
}

type UpdateSensorTypeRequest struct {
	FieldName *string `json:"fieldName"` // 传感器字段名称
}
