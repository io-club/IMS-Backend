package main

import (
	"IMS-Backend/pkgs/fish-net/initialize"
	"IMS-Backend/pkgs/fish-net/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// init
func init() {
	initialize.InitAll()
}

func main() {
	g := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))

	g.Use(middleware.Cors()) // 直接放行全部跨域请求
	// g.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求

	// register router
	register(g)

	g.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
