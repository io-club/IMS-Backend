package api

import (
	"github.com/gin-gonic/gin"
	"ims-server/device/internal/bll"
	ioginx "ims-server/pkg/ginx"
)

func GetDataByID() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewDataService().GetDataByID)
}

// TODO：待做，@姚礼兴
