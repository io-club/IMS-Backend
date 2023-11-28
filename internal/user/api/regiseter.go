package api

import (
	"github.com/gin-gonic/gin"
	"ims-server/internal/user/bll/service"
	ioginx "ims-server/pkg/ginx"
)

func Register() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().Register)
}

func SendVerification() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().SendVerification)
}

func NameLogin() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().NameLogin)
}

func EmailLogin() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().EmailLogin)
}

func RetrievePassword() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().RetrievePassword)
}
