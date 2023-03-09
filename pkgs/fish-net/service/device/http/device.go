package http

import (
	"fishnet/domain"
	"fishnet/glb"
	"fishnet/service/device/pack"
	"fishnet/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateDevice(c *gin.Context) {
	var req pack.CreateDeviceRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	devices, err := _deviceUsecase.QueryDevice(nil, &req.Name, nil, nil, 1, 0)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(devices) > 0 {
		c.JSON(400, gin.H{
			"msg": "device already exists",
		})
		return
	}
	device := &domain.Device{
		GroupID:  0,
		Name:     req.Name,
		Remark:   req.Remark,
		Disabled: false,
	}
	if err := _deviceUsecase.CreateDevice(device); err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"id": device.ID,
		},
	})
}

func DeleteDevice(c *gin.Context) {
	deviceIDStr := c.Param("deviceId")
	if deviceIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	deviceId, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	err = _deviceUsecase.DeleteDevice(deviceId)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func UpdateDevice(c *gin.Context) {
	deviceIDStr := c.Param("deviceId")
	if deviceIDStr == "" {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	deviceId, err := strconv.ParseInt(deviceIDStr, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	var req pack.UpdateDeviceRequest
	if err := c.BindJSON(&req); err != nil {
		glb.LOG.Info("bind json failed")
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	glb.LOG.Info(util.SPrettyLog(req))
	err = _deviceUsecase.UpdateDevice(deviceId, &req.Name, &req.Remark, &req.Disable)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}

func QueryDevice(c *gin.Context) {
	if deviceIDStr := c.Param("deviceId"); deviceIDStr != "" {
		deviceID, err := strconv.ParseInt(deviceIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "bad request",
			})
			return
		}
		devices, err := _deviceUsecase.QueryDevice(&deviceID, nil, nil, nil, 1, 0)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": "internal server error",
			})
			return
		}
		if len(devices) == 0 {
			c.JSON(404, gin.H{
				"msg": "device not found",
			})
			return
		}
		c.JSON(200, gin.H{
			"msg":  "ok",
			"data": pack.Device(devices[0]),
		})
		return
	}

	glb.LOG.Info("aaa multi")
	var req pack.QueryDeviceRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(400, gin.H{
			"msg": "bad request",
		})
		return
	}
	devices, err := _deviceUsecase.QueryDevice(nil, &req.Name, &req.Remark, &req.Disable, req.Limit, req.Offset)
	if err != nil {
		c.JSON(500, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"msg": "ok",
		"data": gin.H{
			"list":  pack.Devices(devices),
			"total": len(devices),
		},
	})
}
