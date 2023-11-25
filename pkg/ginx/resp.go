package ioginx

import (
	"context"
	egoerror "ims-server/pkg/error"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewOk(ctx context.Context, data interface{}) Response {
	return Response{
		Code: egoerror.Ok.Code,
		Data: data,
		Msg:  egoerror.Ok.Zh,
	}
}

func NewErr(ctx context.Context, e egoerror.ErrCode) Response {
	if e == nil {
		return Response{
			Code: egoerror.ErrSystemError.Code,
			Data: nil,
			Msg:  egoerror.ErrSystemError.Zh,
		}
	}
	return Response{
		Code: e.Code,
		Data: nil,
		Msg:  e.Zh,
	}
}

type ListResponse struct {
	Total int64         `json:"total"`
	List  []interface{} `json:"list"`
}
