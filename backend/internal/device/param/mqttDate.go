package param

import (
	"ims-server/internal/device/job"
)

type MqttDateResponse struct {
	Type       job.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string         `json:"TerminalID" form:"TerminalID" binding:"required"`
}

type CreateMqttDateResponse struct {
	MqttDateResponse
}

type GetMqttDateByID struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetMqttDateByIDResponse struct {
	MqttDateResponse
}

type MGetMqttDateByIDsRequest struct {
	IDs []uint `json:"ids" form:"ids" binding:"required"`
}

type MGetMqttDateByIDsResponse struct {
	List []GetMqttDateByIDResponse `json:"list"`
}

type UpdateMqttDateByIDRequest struct {
	ID         uint           `json:"id" form:"id" binding:"required"`
	Type       job.SensorType `json:"type" form:"type"`
	TerminalID string         `json:"TerminalID" form:"TerminalID"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate"`
}

type UpdateMqttDateByIDResponse struct {
	MqttDateResponse
}

type DeleteMqttDateByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteMqttDateByIDResponse struct {
}
