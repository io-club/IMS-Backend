package api

import (
	"github.com/gin-gonic/gin"
	"ims-server/internal/user/bll/service"
	ioginx "ims-server/pkg/ginx"
)

func GetUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().GetUserByID)
}

func MGetUserByIDs() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().MGetUserByIDs)
}

func GetUsers() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().GetUsers)
}

func UpdateUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().UpdateUserByID)
}

func DeleteUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().DeleteUserByID)
}
