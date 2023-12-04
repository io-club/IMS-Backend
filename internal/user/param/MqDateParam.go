package param

import (
	MqCo "ims-server/pkg/Mqtt"
)

type MqDateResponse struct{
	Type       MqCo.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}     `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string          `json:"TerminalID" form:"TerminalID" binding:"required"`
}

