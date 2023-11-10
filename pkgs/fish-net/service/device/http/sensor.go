package http

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/glb"
	"IMS-Backend/pkgs/fish-net/service/device/pack"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateSensor(c *gin.Context) {
	var deviceID int64
	deviceIDStr := c.Param("deviceId")
	glb.LOG.Info("bind json failed [" + deviceIDStr + "]")
	if deviceIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "device id is required",
		})
		return
	}

	deviceIDParsed, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "device id is invalid",
		})
		return
	}
	deviceID = deviceIDParsed
	devices, err := _deviceUsecase.QueryDevice(&deviceID, nil, nil, nil, 1, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(devices) == 0 {
		c.JSON(400, gin.H{
			"msg": "device not found",
		})
		return
	}

	var req pack.CreateSensorRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "body is invalid",
		})
		return
	}
	// TODO check sensor type

	sensor := &domain.Sensor{
		DeviceID:   deviceID,
		WordcaseID: req.WordcaseID,
		Remark:     req.Remark,
		Stat:       0,
	}
	if err := _sensorUsecase.CreateSensor(sensor); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"id": sensor.ID,
		},
	})
}

func DeleteSensor(c *gin.Context) {
	sensorIDStr := c.Param("sensorId")
	if sensorIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	sensorID, err := strconv.ParseInt(sensorIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	if err := _sensorUsecase.DeleteSensor(sensorID); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func UpdateSensor(c *gin.Context) {
	var req pack.UpdateSensorRequest
	sensorIDStr := c.Param("sensorId")
	if sensorIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	sensorID, err := strconv.ParseInt(sensorIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	if err := _sensorUsecase.UpdateSensor(sensorID, &req.Remark, &req.Stat); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func QuerySensor(c *gin.Context) {
	var sensorID int64
	var deviceID int64
	if sensorIDStr := c.Param("sensorId"); sensorIDStr != "" {
		sensorIDParsed, err := strconv.ParseInt(sensorIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "bad request",
			})
			return
		}
		sensorID = sensorIDParsed
	}
	if deviceIDStr := c.Param("deviceId"); deviceIDStr != "" {
		deviceIDParsed, err := strconv.ParseInt(deviceIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "bad request",
			})
			return
		}
		deviceID = deviceIDParsed
	}

	var req pack.QuerySensorRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}

	sensors, err := _sensorUsecase.QuerySensor(&sensorID, &deviceID, &req.SensorType, &req.Stat, 0, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}

	if sensorID != 0 {
		if len(sensors) == 0 {
			c.JSON(400, gin.H{
				"msg": "sensor not found",
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "ok",
			"data": pack.Sensor(sensors[0]),
		})
		return
	}

	if deviceID != 0 {
		c.JSON(200, gin.H{
			"msg": "ok",
			"data": gin.H{
				"list":  pack.Sensors(sensors),
				"total": len(sensors),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"list":  pack.Sensors(sensors),
			"total": len(sensors),
		},
	})
	return

}
