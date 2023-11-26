package api

import (
	"github.com/gin-gonic/gin"
	"ims-server/internal/user/bll/service"
	ioginx "ims-server/pkg/ginx"
)

// TODO: 删除私有创建示例
func CreateUser() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().CreateUser)
}

func GetUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().GetUserByID)
}

func MGetUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().MGetUserByID)
}

func UpdateUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().UpdateUserByID)
}

func DeleteUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(service.NewUserService().DeleteUserByID)
}
