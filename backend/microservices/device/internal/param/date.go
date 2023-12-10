package param

import (
	"ims-server/microservices/device/internal/job"
)

type DataResponse struct {
	SensorType string      `json:"sensorType" form:"sensorType"`
	Date       interface{} `json:"date" form:"date"`
	TerminalID string      `json:"terminalID" form:"terminalID"`
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
	SensorType job.SensorType `json:"sensorType" form:"sensorType"`
	TerminalID string         `json:"terminalID" form:"terminalID"`
	Date       interface{}    `json:"date" form:"date"`
}

type UpdateDataByIDResponse struct {
	DataResponse
}

type DeleteDataByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteDataByIDResponse struct {
}
