package pack

import (
	"ims-server/internal/user/dal/model"
	"ims-server/internal/user/param"
)

func ToMqDateResponse(s *model.Sensor) param.MqDateResponse {
	res := param.MqDateResponse{
		Type: s.Type,
		TerminalID: s.TerminalID,
		SensorDate: s.SensorDate,
	}
	return res
}
