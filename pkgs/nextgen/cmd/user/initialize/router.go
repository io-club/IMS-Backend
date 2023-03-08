package initialize

// func initRouter() {
// 	g := glb.G

// 	g.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "pong",
// 		})
// 	})

// 	v1 := g.Group("/v1")
// 	{
// 		// auth
// 		auth := v1.Group("/auth")
// 		auth.GET("/play", userApi.Play)
// 		auth.GET("/register/begin/:username", userApi.RegisterBegin)
// 		auth.POST("/register/finish/:id", userApi.RegisterFinish)
// 		auth.GET("/login/begin/:username", userApi.LoginBegin)
// 		auth.POST("/login/finish/:id", userApi.LoginFinish)
// 		// user
// 	}

// }
