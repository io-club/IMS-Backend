package pack

import (
	"ims-server/microservices/device/internal/dal/model"
	"ims-server/microservices/device/internal/param"
)

func ToDataResponse(s *model.Data) param.DataResponse {
	res := param.DataResponse{
		SensorType: s.SensorType,
		TerminalID: s.TerminalID,
		Date:       s.Data,
	}
	return res
}
