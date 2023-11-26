package param

import (
	"context"
	ioconsts "ims-server/pkg/consts"
)

type UserResponse struct {
	ID   uint
	Type ioconsts.UserType

	Name     string
	Nickname string
	Avatar   string

	PhoneNumber string
	Email       string

	Status ioconsts.AccountStatus
}

type CreateUserRequest struct {
	Type ioconsts.UserType `json:"type" binding:"required"`

	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`

	Nickname    string `json:"nickname"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
}

type CreateUserResponse struct {
	UserResponse
}

type GetUserByIDRequest struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetUserByIDResponse struct {
	UserResponse
}

type MGetUserByIDRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

type MGetUserByIDResponse struct {
	// TODO: 加上 total?
	List []GetUserByIDResponse `json:"list"`
}

type UpdateUserByIDRequest struct {
	ID   uint              `json:"id" binding:"required"`
	Type ioconsts.UserType `json:"type"`

	Name     string `json:"name"`
	Nickname string `json:"nickname"`

	PhoneNumber string
	Email       string

	Status ioconsts.AccountStatus `json:"accountStatus"`
}

type UpdateUserByIDResponse struct {
	UserResponse
}

// TODO: 更新头像

type DeleteUserByIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

type DeleteUserByIDResponse struct {
}

type IUserService interface {
	// 根据主键 ID 获取用户
	GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*GetUserByIDResponse, error)
	// 根据用户 ID 列表获取多个用户信息
	MGetUserByID(ctx context.Context, req *MGetUserByIDRequest) (*MGetUserByIDResponse, error)
	// 根据主键 ID 更新用户
	UpdateUserByID(ctx context.Context, req *UpdateUserByIDRequest) (*UpdateUserByIDResponse, error)
	// 根据主键 ID 删除用户
	DeleteUserByID(ctx context.Context, req *DeleteUserByIDRequest) (*DeleteUserByIDResponse, error)
}
