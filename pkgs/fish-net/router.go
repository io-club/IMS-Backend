package main

import (
	confAPI "fishnet/service/conf/http"
	deviceAPI "fishnet/service/device/http"
	userAPI "fishnet/service/user/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

func register(g *gin.Engine) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	type User struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	v1 := g.Group("/v1")
	{
		// auth
		auth := v1.Group("/auth")
		{
			auth.GET("/play", userAPI.Play)
			auth.GET("/register/begin/:username", userAPI.RegisterBegin)
			auth.POST("/register/finish/:id", userAPI.RegisterFinish)
			auth.GET("/login/begin/:username", userAPI.LoginBegin)
			auth.POST("/login/finish/:id", userAPI.LoginFinish)
		}

		// user
		user := v1.Group("/user")
		{
			user.POST("", userAPI.CreateUser)
			user.DELETE("/:userId", userAPI.DeleteUser)
			user.PUT("/:userId", userAPI.UpdateUser)
			user.GET("", userAPI.QueryUser)
			user.GET("/:userId", userAPI.QueryUser)
		}

		// wordcase
		wordcase := v1.Group("/wordcase")
		{
			wordcase.POST("", confAPI.CreateWordcase)
			wordcase.DELETE("/:wordcaseId", confAPI.DeleteWordcase)
			wordcase.PUT("/:wordcaseId", confAPI.UpdateWordcase)
			wordcase.GET("", confAPI.QueryWordcase)
			wordcase.GET("/:wordcaseId", confAPI.QueryWordcase)
		}

		// device
		device := v1.Group("/device")
		{
			// device
			device.POST("", deviceAPI.CreateDevice)
			device.DELETE("/:deviceId", deviceAPI.DeleteDevice)
			device.PUT("/:deviceId", deviceAPI.UpdateDevice)
			device.GET("", deviceAPI.QueryDevice)
			device.GET("/:deviceId", deviceAPI.QueryDevice)

			deviceSensor := device.Group("/:deviceId/sensor")
			{
				deviceSensor.POST("", deviceAPI.CreateSensor)
				deviceSensor.GET("", deviceAPI.QuerySensor)
				deviceSensor.GET("/:sensorId", deviceAPI.QuerySensor)
			}
		}

		// sensor
		sensor := v1.Group("/sensor")
		{
			sensor.DELETE("/:sensorId", deviceAPI.DeleteSensor)
			sensor.PUT("/:sensorId", deviceAPI.UpdateSensor)
			sensor.GET("", deviceAPI.QuerySensor)
			sensor.GET("/:sensorId", deviceAPI.QuerySensor)
		}

	}
}
