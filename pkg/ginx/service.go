package ioginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	ioerror "ims-server/pkg/error"
	iologger "ims-server/pkg/logger"
	"net/http"
	"reflect"
)

func BindRequest(c *gin.Context, svc interface{}) (interface{}, ioerror.ErrCode) {
	svcType := reflect.TypeOf(svc)
	if svcType.Kind() != reflect.Func {
		iologger.Info("svc is not a function")
		return nil, ioerror.ErrSystemError
	}
	if svcType.NumIn() != 2 {
		iologger.Info("svc parameter count is not 2")
		return nil, ioerror.ErrSystemError
	}
	// If it is a `GET` request, only use the `Form` binding engine (`query`).
	// If it is a `POST` request, first check if `content-type` is `JSON` or `XML`, and then use `Form` (`form-data`).
	reqType := svcType.In(1).Elem()
	req := reflect.New(reqType)
	if err := c.ShouldBind(req.Interface()); err != nil {
		iologger.Info("bind Params failed,err: %v", err)
		return nil, ioerror.ErrInvalidParam
	}
	return req.Interface(), nil
}

func CallaService(ctx context.Context, req interface{}, svc interface{}) ([]reflect.Value, ioerror.ErrCode) {
	if reflect.TypeOf(svc).Kind() != reflect.Func {
		iologger.Info("svc is not a function")
		return nil, ioerror.ErrSystemError
	}
	params := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(req)}
	rets := reflect.ValueOf(svc).Call(params) // Call svc with params as input and return rets
	if len(rets) != 2 {
		iologger.Info("svc return parameter count is not 2")
		return nil, ioerror.ErrSystemError
	}
	return rets, nil
}

func CheckResponse(rets []reflect.Value) (interface{}, ioerror.ErrCode) {
	// The service function must return two values
	if len(rets) != 2 {
		iologger.Info("service return parameter count is not 2")
		return nil, ioerror.ErrSystemError
	}

	// If error is nil, return directly
	if rets[1].IsNil() {
		return rets[0].Interface(), nil
	}

	// If it is an ErrCode type, return directly
	if err, ok := rets[1].Interface().(ioerror.ErrCode); ok {
		if err == nil || err.Code == 0 {
			return rets[0].Interface(), nil
		}
		return nil, err
	}

	// If it is another error type, return system error
	iologger.Warn("service return error is not ErrCode")
	return nil, ioerror.ErrSystemError
}

func ToHandle(fn interface{}) func(c *gin.Context) {
	handler := func(c *gin.Context) {
		if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
			err := v.RegisterValidation("password", CheckPassword)
			if err != nil {
				iologger.Debug("validator register failed,err: %v", err)
				return
			}
		}

		ctx := context.Background()
		// Parse request parameters
		req, err := BindRequest(c, fn)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErr(c, err))
			return
		}
		iologger.Info("request: %#v", req)

		// Call the handler function
		rets, err := CallaService(ctx, req, fn)
		if err != nil {
			iologger.Debug("CallaService Failed")
			c.JSON(http.StatusBadRequest, NewErr(c, err))
			return
		}

		// Process the response
		res, err := CheckResponse(rets)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErr(c, err))
			return
		}
		iologger.Info("response: %#v", res)
		c.JSON(http.StatusOK, NewOk(c, res))
	}
	return handler
}
