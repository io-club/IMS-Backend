package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	ioerror "ims-server/pkg/error"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/util"
	"ims-server/user/internal/bll"
	"ims-server/user/internal/param"
	"net/http"
	"time"
)

func Register() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().Register)
}

func SendVerification() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().SendVerification)
}

func NameLogin() func(ctx *gin.Context) {
	handler := func(c *gin.Context) {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation("password", ioginx.CheckPassword)
			if err != nil {
				iologger.Debug("validator register failed,err: %v", err)
				c.JSON(http.StatusBadRequest, ioginx.NewErr(c, ioerror.ErrInvalidParam))
				return
			}
		}
		// 绑定参数
		var req param.NameLoginRequest
		err := c.ShouldBind(&req)
		if err != nil {
			iologger.Debug("bind Params failed,err: %v", err)
			c.JSON(http.StatusBadRequest, ioginx.NewErr(c, ioerror.ErrInvalidParam))
			return
		}
		// 处理请求
		resp, err := bll.NewUserService().NameLogin(c, &req)
		if err != nil {
			var ioErr ioerror.ErrCode
			ok := errors.As(err, &ioErr)
			if !ok {
				ioErr = ioerror.ErrSystemError
			}
			c.JSON(http.StatusBadRequest, ioErr)
			return
		}
		// 生成 refresh token
		payload := util.JwtPayload{
			Issue:       "ims",
			Audience:    "ims",
			Subject:     "refresh-Token",
			IssueAt:     time.Now().Unix(),
			Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(),
			UserDefined: map[string]any{"uid": resp.ID, "uname": resp.Name, "utype": resp.Type},
		}
		refreshToken, err := util.GenJwt(util.DefautHeader, payload)
		// 生成 access token
		payload.Subject = "access-Token"
		payload.Expiration = time.Now().Add(1 * time.Hour).Unix()
		accessToken, err := util.GenJwt(util.DefautHeader, payload)

		c.Header("access-Token", accessToken)
		c.Header("refresh-Token", refreshToken)
		// 生成响应
		c.JSON(http.StatusOK, ioginx.NewOk(c, resp))
	}
	return handler
}

func EmailLogin() func(ctx *gin.Context) {
	handler := func(c *gin.Context) {
		// 绑定参数
		var req param.EmailLoginRequest
		err := c.ShouldBind(&req)
		if err != nil {
			iologger.Debug("bind Params failed,err: %v", err)
			c.JSON(http.StatusBadRequest, ioginx.NewErr(c, ioerror.ErrInvalidParam))
			return
		}

		// 处理请求
		resp, err := bll.NewUserService().EmailLogin(c, &req)
		if err != nil {
			var ioErr ioerror.ErrCode
			ok := errors.As(err, &ioErr)
			if !ok {
				ioErr = ioerror.ErrSystemError
			}
			c.JSON(http.StatusBadRequest, ioginx.NewErr(c, ioErr))
			return
		}
		// 生成 refresh token
		payload := util.JwtPayload{
			Issue:       "ims",
			Audience:    "ims",
			Subject:     "refresh-Token",
			IssueAt:     time.Now().Unix(),
			Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(),
			UserDefined: map[string]any{"uid": resp.ID, "uname": resp.Name, "utype": resp.Type},
		}
		refreshToken, err := util.GenJwt(util.DefautHeader, payload)
		// 生成 access token
		payload.Subject = "access-Token"
		payload.Expiration = time.Now().Add(1 * time.Hour).Unix()
		accessToken, err := util.GenJwt(util.DefautHeader, payload)

		c.Header("access-Token", accessToken)
		c.Header("refresh-Token", refreshToken)
		// 生成响应
		c.JSON(http.StatusOK, ioginx.NewOk(c, resp))
	}
	return handler
}

func RetrievePassword() func(ctx *gin.Context) {
	return ioginx.ToHandle(bll.NewUserService().RetrievePassword)
}
