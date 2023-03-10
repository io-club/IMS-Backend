package http

import (
	"fishnet/domain"
	"fishnet/glb"
	"fishnet/service/device/pack"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateSensorData(c *gin.Context) {
	var req pack.CreateSensorDataRequest
	if err := c.BindJSON(&req); err != nil {
		glb.LOG.Info("bind json failed [" + err.Error() + "]")
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	if req.SensorID == nil || *req.SensorID == 0 || req.Value == "" || req.CollectTime.IsZero() {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	// check sensor exists
	sensors, err := _sensorUsecase.QuerySensor(
		req.SensorID,
		nil, nil, nil, 1, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(sensors) == 0 {
		c.JSON(400, gin.H{
			"msg": "sensor not exists",
		})
		return
	}
	// check sensor data has been created
	sensorDatas, err := _sensorDataUsecase.QuerySensorData(nil, req.SensorID, req.SensorTypeID, nil, &req.CollectTime, nil, nil, 1, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(sensorDatas) > 0 {
		c.JSON(400, gin.H{
			"msg": "sensor data has been created",
		})
		return
	}
	// check sensor type exists
	sensorTypes, err := _sensorTypeUsecase.QuerySensorType(req.SensorTypeID, nil, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(sensorTypes) == 0 {
		c.JSON(400, gin.H{
			"msg": "sensor type not exists",
		})
		return
	}

	sensorData := &domain.SensorData{
		SensorID:     *req.SensorID,
		SensorTypeID: *req.SensorTypeID,
		Value:        req.Value,
		CollectAt:    req.CollectTime,
	}
	if err := _sensorDataUsecase.CreateSensorData([]*domain.SensorData{sensorData}); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"id": sensorData.ID,
		},
	})
}

func DeleteSensorData(c *gin.Context) {
	sensorDataIDStr := c.Param("sensorDataId")
	if sensorDataIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	sensorDataID, err := strconv.ParseInt(sensorDataIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	if err := _sensorDataUsecase.DeleteSensorData([]int64{sensorDataID}); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func UpdateSensorData(c *gin.Context) {
	sensorDataIDStr := c.Param("sensorDataId")
	if sensorDataIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	sensorDataID, err := strconv.ParseInt(sensorDataIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	var req pack.UpdateSensorDataRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	if err := _sensorDataUsecase.UpdateSensorData(sensorDataID, nil, nil, req.Value); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func QuerySensorData(c *gin.Context) {
	sensorDataIDStr := c.Param("sensorDataId")
	var sensorDataID *int64
	if sensorDataIDStr != "" {
		sensorDataIDInt, err := strconv.ParseInt(sensorDataIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "bad request",
			})
			return
		}
		if sensorDataIDInt != 0 {
			sensorDataID = &sensorDataIDInt
		} else {
			sensorDataID = nil
		}
	}
	var req pack.QuerySensorDataRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	// Query sensor type
	sensorTypes, err := _sensorTypeUsecase.QuerySensorType(nil, req.WordcaseID, nil)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	sensorTypeID2SensorType := make(map[int64]*domain.SensorType)
	for _, sensorType := range sensorTypes {
		sensorTypeID2SensorType[int64(sensorType.ID)] = sensorType
	}
	// Query sensor data
	sensorDatas, err := _sensorDataUsecase.QuerySensorData(
		sensorDataID,
		req.SensorId,
		req.WordcaseID,
		nil,
		nil, req.BeginTime, req.EndTime, req.Limit, req.Offset)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	// Pack response
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"list":  pack.SensorDatas(sensorDatas, sensorTypeID2SensorType),
			"total": len(sensorDatas),
		},
	})
}
