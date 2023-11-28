package param

import (
	"context"
	ioconsts "ims-server/pkg/consts"
)

type UserResponse struct {
	ID   uint              `json:"id" form:"id"`
	Type ioconsts.UserType `json:"type" form:"type"`

	Name     string `json:"name" form:"name"`
	Nickname string `json:"nickname" form:"nickname"`
	Avatar   string `json:"avatar" form:"avatar"`

	PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`
	Email       string `json:"email" form:"email"`

	Status ioconsts.AccountStatus `json:"status" form:"status"`
}

type GetUserByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetUserByIDResponse struct {
	UserResponse
}

type MGetUserByIDsRequest struct {
	IDs []uint `json:"ids" form:"ids" binding:"required"`
}

type MGetUserByIDsResponse struct {
	// TODO: 加上 total?
	List []GetUserByIDResponse `json:"list"`
}

type UpdateUserByIDRequest struct {
	ID   uint              `json:"id" form:"id" binding:"required"`
	Type ioconsts.UserType `json:"type" form:"type"`

	Name     string `json:"name" form:"name"`
	Nickname string `json:"nickname" form:"nickname"`

	PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`
	Email       string `json:"email" form:"email"`

	Status ioconsts.AccountStatus `json:"accountStatus" form:"accountStatus"`
}

type UpdateUserByIDResponse struct {
	UserResponse
}

// TODO: 更新头像

type DeleteUserByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type DeleteUserByIDResponse struct {
}

type IUserService interface {
	// 根据主键 ID 获取用户
	GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*GetUserByIDResponse, error)
	// 根据用户 ID 列表获取多个用户信息
	MGetUserByIDs(ctx context.Context, req *MGetUserByIDsRequest) (*MGetUserByIDsResponse, error)
	// 根据主键 ID 更新用户
	UpdateUserByID(ctx context.Context, req *UpdateUserByIDRequest) (*UpdateUserByIDResponse, error)
	// 根据主键 ID 删除用户
	DeleteUserByID(ctx context.Context, req *DeleteUserByIDRequest) (*DeleteUserByIDResponse, error)

	// 注册
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	// 发送验证码
	SendVerification(ctx context.Context, req *SendVerification) (*SendVerificationResponse, error)
}
