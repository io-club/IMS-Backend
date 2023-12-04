package param

import (
	"ims-server/internal/device/job"
)

type MqDateResponse struct {
	Type       job.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string         `json:"TerminalID" form:"TerminalID" binding:"required"`
}
