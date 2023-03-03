package initialize

import (
	"demo/glb"
	"net/http"

	"demo/api"

	"github.com/gin-gonic/gin"
)

func initRouter() {
	glb.G.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	g := glb.G
	v1 := g.Group("/auth")
	{
		v1.GET("/register/play", api.Play)
		v1.GET("/register/begin/:username", api.RegisterBegin)
		v1.POST("/register/finish/:id", api.RegisterFinish)
		// v1.PUT("/:id", updateTodo)
		// v1.DELETE("/:id", deleteTodo)
	}

}
