package ioginx

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	ioerror "ims-server/pkg/error"
	iologger "ims-server/pkg/logger"
	"ims-server/pkg/registery"
	"ims-server/pkg/util"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
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

		// Parse request parameters
		req, err := BindRequest(c, fn)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewErr(c, err))
			return
		}
		iologger.Info("request: %#v", req)

		// Call the handler function
		rets, err := CallaService(c, req, fn)
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

func post(ctx context.Context, url string, data interface{}, headers map[string]string) (*Response, error) {
	restyRes, err := resty.New().R().
		SetHeaders(headers).
		SetContext(ctx).
		SetBody(data).
		SetResult(&Response{}).
		Post(url)
	if err != nil {
		return nil, err
	}
	statusCode := restyRes.StatusCode()
	if statusCode != 200 {
		return nil, errors.New("response status code is not 200")
	}
	if res, ok := restyRes.Result().(*Response); ok {
		return res, nil
	}
	return nil, errors.New("failed to parse response")
}

func Call[T interface{}](ctx context.Context, serviceName string, req T) (*T, error) {
	// Service discovery
	hub := registery.GetServiceClient()
	servers := hub.GetServiceEndpointsWithCache(registery.UserService)
	// TODO: Implement a more efficient load balancing algorithm
	// Load balancing: select a service node
	idx := rand.Intn(len(servers))
	addr := servers[idx]
	// Extract function name
	reqType := reflect.TypeOf(req)
	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
	}
	reqName := reqType.Name()
	reqName = strings.ToLower(reqName)
	fn := strings.TrimSuffix(reqName, "request")
	// Build URL
	url := fmt.Sprintf("http://%s/%s", addr, fn)
	// Build headers
	payload := util.JwtPayload{
		Issue:       "ims",
		Audience:    "ims",
		Subject:     "refresh-Token",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(),
		UserDefined: map[string]any{"uid": ctx.Value("uid"), "uname": ctx.Value("uname"), "utype": ctx.Value("utype")},
	}
	accessToken, err := util.GenJwt(util.DefautHeader, payload)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}
	headers := map[string]string{
		"access-Token": accessToken,
	}
	iologger.Info("call url: %s", url)
	resp, err := post(ctx, url, req, headers)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		iologger.Debug("error returned from remote interface: %s", resp.Msg)
		return nil, errors.New("error returned from remote interface")
	}
	s, err := sonic.Marshal(resp.Data)
	if err != nil {
		iologger.Warn("failed to serialize the returned data, err: %v", err)
		return nil, errors.New("failed to serialize the returned data")
	}
	var res T
	err = sonic.Unmarshal(s, &res)
	if err != nil {
		iologger.Warn("failed to deserialize the returned data, err: %v", err)
		return nil, errors.New("failed to deserialize the returned data")
	}
	return &res, nil
}
