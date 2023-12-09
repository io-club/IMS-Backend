package pack

import (
	"ims-server/microservices/device/internal/dal/model"
	"ims-server/microservices/device/internal/param"
)

func ToDataResponse(s *model.Data) param.DataResponse {
	res := param.DataResponse{
		Type:       s.Type,
		TerminalID: s.TerminalID,
		SensorDate: s.SensorData,
	}
	return res
}
