package ioginx

import (
	"context"
	"github.com/gin-gonic/gin"
	ioerror "ims-server/pkg/error"
	iologger "ims-server/pkg/logger"
	"reflect"
)

func BindRequest(c *gin.Context, route Route) (interface{}, ioerror.ErrCode) {
	svc := route.Func
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
		iologger.Info("bind Params failed")
		return nil, ioerror.ErrInvalidParam
	}

	return req.Interface(), nil
}

func ToHandle(ctx *context.Context, req interface{}, svc interface{}) ([]reflect.Value, ioerror.ErrCode) {
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

func BindResponse(rets []reflect.Value) (interface{}, ioerror.ErrCode) {
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
