package param

import (
	"ims-server/device/internal/job"
)

type DataResponse struct {
	Type       job.SensorType `json:"type" form:"type" binding:"required"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate" binding:"required"`
	TerminalID string         `json:"TerminalID" form:"TerminalID" binding:"required"`
}

type CreateDataResponse struct {
	DataResponse
}

type GetDataByID struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetDataByIDResponse struct {
	DataResponse
}

type MGetDataByIDsRequest struct {
	IDs []uint `json:"ids" form:"ids" binding:"required"`
}

type MGetDataByIDsResponse struct {
	List []GetDataByIDResponse `json:"list"`
}

type UpdateDataByIDRequest struct {
	ID         uint           `json:"id" form:"id" binding:"required"`
	Type       job.SensorType `json:"type" form:"type"`
	TerminalID string         `json:"TerminalID" form:"TerminalID"`
	SensorDate interface{}    `json:"SensorDate" form:"SensorDate"`
}

type UpdateDataByIDResponse struct {
	DataResponse
}

type DeleteDataByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteDataByIDResponse struct {
}
