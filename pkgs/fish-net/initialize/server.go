package initialize

import (
	"fishnet/glb"

	"fishnet/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func initServer() {
	g := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))

	g.Use(middleware.Cors()) // 直接放行全部跨域请求
	// g.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	glb.G = g
}
