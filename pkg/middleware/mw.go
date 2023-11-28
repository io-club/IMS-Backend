package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	iologger "ims-server/pkg/logger"
	"time"
)

// Limit the maximum concurrency of each interface
var limitCh = make(chan struct{}, 500)

func TimeMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		ctx.Next()
		timeCost := time.Since(now)
		// TODO: 信息收集
		msg := fmt.Sprintf("request %s cost time %d ms\n", ctx.Request.URL.Path, timeCost.Milliseconds())
		iologger.Info(msg)
	}
}

func LimitMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limitCh <- struct{}{} // When the maximum concurrency limit is reached, block
		ctx.Next()
		<-limitCh
	}
}
