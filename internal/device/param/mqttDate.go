package param

import (
	"ims-server/internal/device/job"
)

type MqDateRequest struct {
	Type       job.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string         `json:"TerminalID" form:"TerminalID" binding:"required"`
}

type MqttDateResponse struct {
	Type       job.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string         `json:"TerminalID" form:"TerminalID" binding:"required"`
}

type CreateMqttDateRes struct {
	MqttDateResponse
}

type GetMqttDateByID struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetMqttDateByIDRes struct {
	MqttDateResponse
}

type MGetMqttDateByIDsReq struct {
	IDs []uint `json:"ids" form:"ids" binding:"required"`
}

type MGetMqttDateByIDsRes struct {
	List []GetMqttDateByIDRes `json:"list"`
}

type UpdateMqttDateByIDReq struct {
	ID         uint           `json:"id" form:"id" binding:"required"`
	Type       job.SensorType `json:"type" form:"type"`
	TerminalID string         `json:"TerminalID" form:"TerminalID"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate"`
}

type UpdateMqttDateByIDRes struct {
	MqttDateResponse
}

type DeleteMqttDateByIDReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteMqttDateByIDRes struct {
}
