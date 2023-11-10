package http

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/device/pack"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateSensorType(c *gin.Context) {
	var req pack.CreateSensorTypeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "bad request",
		})
		return
	}
	wordcases, err := _wordcaseUsecase.QueryWordcase(&req.WordcaseID, nil, nil, 1, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(wordcases) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "sensor type not exists",
		})
		return
	}

	if len(req.FieldName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "field name is empty",
		})
		return
	}

	wordcaseID := int64(wordcases[0].ID)
	sensorTypes, err := _sensorTypeUsecase.QuerySensorType(nil, &wordcaseID, &req.FieldName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(sensorTypes) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "field already exists",
		})
		return
	}
	sensorType := &domain.SensorType{
		WordcaseID: req.WordcaseID,
		FieldName:  req.FieldName,
		FieldType:  req.FieldType,
	}
	if err := _sensorTypeUsecase.CreateSensorType(sensorType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"data": gin.H{
			"id": sensorType.ID,
		},
	})
}

func DeleteSensorType(c *gin.Context) {
	sensorTypeIDStr := c.Param("sensorTypeId")
	sensorTypeID, err := strconv.ParseInt(sensorTypeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "bad request",
		})
		return
	}
	if err := _sensorTypeUsecase.DeleteSensorType(sensorTypeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func UpdateSensorType(c *gin.Context) {
	sensorTypeIDStr := c.Param("sensorTypeId")
	sensorTypeID, err := strconv.ParseInt(sensorTypeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "bad request",
		})
		return
	}
	sensorTypes, err := _sensorTypeUsecase.MGetSensorTypes([]int64{sensorTypeID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(sensorTypes) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "sensor type not exists",
		})
		return
	}
	sensorType := sensorTypes[0]
	var req pack.UpdateSensorTypeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "bad request",
		})
		return
	}
	if req.FieldName == nil || len(*req.FieldName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "field name is empty",
		})
		return
	}
	sensorTypes, err = _sensorTypeUsecase.QuerySensorType(nil, &sensorType.WordcaseID, req.FieldName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	if len(sensorTypes) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "sensor type already exists",
		})
		return
	}
	if err := _sensorTypeUsecase.UpdateSensorType(sensorType.WordcaseID, req.FieldName, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func QuerySensorType(c *gin.Context) {
	// findByID
	sensorTypeIDStr := c.Param("sensorTypeId")
	var sensorTypeID *int64 = nil
	if sensorTypeIDStr != "" {
		sensorTypeIDParsed, err := strconv.ParseInt(sensorTypeIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "bad request",
			})
			return
		}
		sensorTypeID = &sensorTypeIDParsed
	}
	// var req pack.QuerySensorTypeRequest
	wordcaseIDStr := c.Query("wordcaseId")
	var wordcaseID *int64 = nil
	if wordcaseIDStr != "" {
		wordcaseIDParsed, err := strconv.ParseInt(wordcaseIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "bad request",
			})
			return
		}
		wordcaseID = &wordcaseIDParsed
	}

	// Query
	fieldName := c.Query("fieldName")
	sensorTypes, err := _sensorTypeUsecase.QuerySensorType(sensorTypeID, wordcaseID, &fieldName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}

	if sensorTypeID != nil {
		if len(sensorTypes) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "sensor type not exists",
			})
			return
		}
		if len(sensorTypes) > 0 {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"data": pack.SensorType(sensorTypes[0]),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
		"data": gin.H{
			"list":  pack.SensorTypes(sensorTypes),
			"total": len(sensorTypes),
		},
	})
}
