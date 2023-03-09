package main

import (
	confApi "fishnet/service/conf/http"
	deviceApi "fishnet/service/device/http"
	userApi "fishnet/service/user/http"
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
			auth.GET("/play", userApi.Play)
			auth.GET("/register/begin/:username", userApi.RegisterBegin)
			auth.POST("/register/finish/:id", userApi.RegisterFinish)
			auth.GET("/login/begin/:username", userApi.LoginBegin)
			auth.POST("/login/finish/:id", userApi.LoginFinish)
		}

		// user
		user := v1.Group("/user")
		{
			user.POST("", userApi.CreateUser)
			user.DELETE("/:userId", userApi.DeleteUser)
			user.PUT("/:userId", userApi.UpdateUser)
			user.GET("", userApi.QueryUser)
			user.GET("/:userId", userApi.QueryUser)
		}

		// wordcase
		wordcase := v1.Group("/wordcase")
		{
			wordcase.POST("", confApi.CreateWordcase)
			wordcase.DELETE("/:wordcaseId", confApi.DeleteWordcase)
			wordcase.PUT("/:wordcaseId", confApi.UpdateWordcase)
			wordcase.GET("", confApi.QueryWordcase)
			wordcase.GET("/:wordcaseId", confApi.QueryWordcase)
		}

		// device
		device := v1.Group("/device")
		{
			device.POST("", deviceApi.CreateDevice)
			device.DELETE("/:deviceId", deviceApi.DeleteDevice)
			device.PUT("/:deviceId", deviceApi.UpdateDevice)
			device.GET("", deviceApi.QueryDevice)
			device.GET("/:deviceId", deviceApi.QueryDevice)
		}
	}
}
