package pack

import (
	"ims-server/internal/device/dal/model"
	"ims-server/internal/device/param"
)

func ToMqttDateResponse(s *model.SensorData) param.MqttDateResponse {
	res := param.MqttDateResponse{
		Type:       s.Type,
		TerminalID: s.TerminalID,
		SensorDate: s.SensorData,
	}
	return res
}
