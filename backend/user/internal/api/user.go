package api

import (
	"github.com/gin-gonic/gin"
	ioginx "ims-server/pkg/ginx"
	"ims-server/user/internal/bll"
)

func GetUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().GetUserByID)
}

func MGetUserByIDs() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().MGetUserByIDs)
}

func GetUsers() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().GetUsers)
}

func UpdateUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().UpdateUserByID)
}

func UploadAvatar() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().UploadAvatar)
}

func DeleteUserByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().DeleteUserByID)
}
