package ioginx

import (
	"context"
	ioerror "ims-server/pkg/error"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func NewOk(ctx context.Context, data interface{}) Response {
	return Response{
		Code: ioerror.OK.Code,
		Data: data,
		Msg:  ioerror.OK.Zh,
	}
}

func NewErr(ctx context.Context, e ioerror.ErrCode) Response {
	if e == nil {
		return Response{
			Code: ioerror.ErrSystemError.Code,
			Data: nil,
			Msg:  ioerror.ErrSystemError.Zh,
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
