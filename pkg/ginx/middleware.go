package ioginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/util"
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

func JwtAuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("access-Token")
		_, payload, err := util.VerifyJwt(accessToken)
		if err != nil {
			iologger.Debug("verify accessToken failed, err: %v", err)
			c.Next()
			return
		}
		// Check if the token has expired
		if payload.Expiration >= time.Now().Unix() {
			// If not expired, extract custom key-value pairs from the payload
			for k, v := range payload.UserDefined {
				c.Set(k, v)
			}
			c.Next()
		}
		// If expired, try refreshing the accessToken
		refreshToken := c.GetHeader("refresh-Token")
		_, payload, err = util.VerifyJwt(refreshToken)
		if err != nil {
			iologger.Debug("verify refreshToken failed, err: %v", err)
			c.Next()
			return
		}
		// If refreshToken is not expired, generate a new accessToken
		if payload.Expiration >= time.Now().Unix() {
			payload.Subject = "access-Token"
			payload.Expiration = time.Now().Add(1 * time.Hour).Unix()
			accessToken, err = util.GenJwt(util.DefautHeader, *payload)
			if err != nil {
				return
			}
			c.Header("access-Token", accessToken)
			// Also retrieve the custom key-value pairs
			for k, v := range payload.UserDefined {
				c.Set(k, v)
			}
		}
		iologger.Debug("refreshToken has expired")
		c.Next()
		return
	}
}
