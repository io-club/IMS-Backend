package main

import (
	
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type User struct{
	Username string `json:"name" binding:"required"`
	Age int	`json:"age" binding:"min=8,max=99"`
	Password string `json:"password" binding:"required"`
}

func main(){
	r:=gin.Default()
	r.GET("/get/:name",func(c *gin.Context){
		name:=c.Param("name")
		c.JSON(http.StatusOK,gin.H{
			"username":name,
		})
	})
	r.POST("/post",func(c *gin.Context){
		var user User
		if err:=c.ShouldBindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":err.Error(),
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
			"message":"create user successfully",
			})
		}

		
	})
	r.PUT("/update",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"message":"update successful",
		})

	})

	r.DELETE("/delete",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"message":"delete successful",
		})
	})

	r.Run()
}