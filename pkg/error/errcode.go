package ioerror

import (
	"fmt"
)

type ErrCode = *errCode

type errCode struct {
	Code int    `json:"code"`
	Zh   string `json:"zh"`
	En   string `json:"en"`
}

func (e ErrCode) Error() string {
	return fmt.Sprintf("code: %d, msg: %s, advice: %s", e.Code, e.Zh, e.En)
}

func newErrCode(code int, en, zh string) ErrCode {
	return &errCode{
		Code: code,
		Zh:   zh,
		En:   en,
	}
}

// The format of the error code is 000 000. The first three digits represent the module, and the last three digits represent the error code.
// The module number for general error codes is 0.
// For example: |      0|, represents the first error code of the general module, OK.
// For example: |  1 001|, represents the first error code of the user module.
// For example: |111 001|, represents the first error code of the 111th module.

// ---------------------------- 通用错误码 ----------------------------
var (
	ErrQueueCallFailed = newErrCode(-4, "message queue call failed", "消息队列调用失败")
	ErrFfmpegError     = newErrCode(-3, "ffmpeg call failed", "ffmpeg 调用失败")
	ErrDBError         = newErrCode(-3, "db call failed", "数据库调用失败")
	ErrUnimplemented   = newErrCode(-2, "not implemented", "接口未实现")
	ErrSystemError     = newErrCode(-1, "system error", "系统繁忙，此时请稍候再试")

	Ok = newErrCode(0, "success", "成功")

	ErrUnauthorized  = newErrCode(1, "unauthorized", "未登录")
	ErrNotPermitted  = newErrCode(2, "not permitted", "无权限访问")
	ErrInvalidParam  = newErrCode(3, "invalid param", "参数错误")
	ErrNotFound      = newErrCode(4, "not found", "资源不存在")
	ErrRepeatedEntry = newErrCode(5, "repeated entry", "重复录入")
	ErrNoTopic       = newErrCode(6, "topic does not exist", "主题不存在")
	ErrInvalidType   = newErrCode(7, "invalid type", "类型错误")
)

// ---------------------------- 用户模块错误码 ----------------------------

var (
	ErrUserLoginFailed     = newErrCode(1000+1, "login failed", "登录失败")
	ErrInvalidVerifyCode   = newErrCode(1000+2, "invalid verify code", "验证码错误")
	ErrUserNotFound        = newErrCode(1000+3, "user not found", "用户不存在")
	ErrPhoneNumberNotFound = newErrCode(1000+4, "phone number not found", "用户未绑定手机号")
)

// ---------------------------- 设备管理模块错误码 ----------------------------

func NewErrCode(code int) ErrCode {
	// TODO:待做
	switch code {
	// system
	case ErrQueueCallFailed.Code:
		return ErrQueueCallFailed
	case ErrFfmpegError.Code:
		return ErrFfmpegError
	case ErrDBError.Code:
		return ErrDBError
	case ErrUnimplemented.Code:
		return ErrUnimplemented
	case ErrSystemError.Code:
		return ErrSystemError
	case Ok.Code:
		return Ok
	// user
	case ErrUserLoginFailed.Code:
		return ErrUserLoginFailed
	case ErrInvalidVerifyCode.Code:
		return ErrInvalidVerifyCode
	case ErrUserNotFound.Code:
		return ErrUserNotFound
	default:
		return ErrSystemError
	}
}
