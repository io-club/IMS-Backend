package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Login struct {
	// 结构字段必须以大写字母（导出）开头，JSON 包才能看到其值。
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Id       string `form:"id" json:"id" binding:"-"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func main() {
	r := gin.Default()

	// 获取path参数
	r.GET("/:User/:Id", func(c *gin.Context) {
		User := c.Param("User")
		Id := c.Param("Id")
		c.JSON(http.StatusOK, gin.H{
			"User": User,
			"Id":   Id,
		})
	})

	// 获取querystring参数
	// r.GET("/loginForm", func(c *gin.Context) {
	// 	var login Login
	// 	if err := c.ShouldBind(&login); err == nil {
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"user":     login.User,
	// 			"id":       login.Id,
	// 			"password": login.Password,
	// 		})
	// 	} else {
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": err.Error()})
	// 	}
	// })

	//检验并绑定JSON的示例
	r.POST("/loginJson", func(c *gin.Context) {
		var login Login

		if err := c.ShouldBind(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if login.User != "xxx" || login.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":     login.User,
			"id":       login.Id,
			"password": login.Password,
			"Status":   "you are logged in",
		})

	})

	// 检验并绑定form
	r.POST("LoginForm", func(c *gin.Context) {
		var login Login

		if err := c.ShouldBind(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if login.User != "xxx" || login.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":     login.User,
			"id":       login.Id,
			"password": login.Password,
			"Status":   "you are logged in",
		})

	})

	r.PUT("/user", func(c *gin.Context) {
		user := c.DefaultPostForm("user", "Tom")
		id := c.PostForm("321")

		c.JSON(http.StatusOK, gin.H{
			"user": user,
			"id":   id,
		})
	})

	r.DELETE("user/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})
	r.Run(":8080")
}
