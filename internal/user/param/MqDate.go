package param

import (
	MqCo "ims-server/pkg/Mqtt"
)

type MqDateRequest struct {
	Type       MqCo.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}     `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string          `json:"TerminalID" form:"TerminalID" binding:"required"`
}

type CreateMqDateRes struct {
	MqDateResponse
}

type GetMqDateByID struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetMqDateByIDRes struct {
	MqDateResponse
}

type MGetMqDateByIDsReq struct {
	IDs []uint `json:"ids" form:"ids" binding:"required"`
}

type MGetMqDateByIDsRes struct {
	List []GetMqDateByIDRes `json:"list"`
}

type UpdateMqDateByIDReq struct {
	ID         uint            `json:"id" form:"id" binding:"required"`
	Type       MqCo.SensorType `json:"type" form:"type"`
	TerminalID string          `json:"TerminalID" form:"TerminalID"`
	SensorDate interface{}     `json:"SensorDate" form:"SensorDate"`
}

type UpdateMqDateByIDRes struct {
	MqDateResponse
}

type DeleteMqDateByIDReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteMqDateByIDRes struct {
}

