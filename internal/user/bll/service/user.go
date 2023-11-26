package userservice

import (
	"context"
	"ims-server/internal/user/bll/pack"
	"ims-server/internal/user/dal/model"
	"ims-server/internal/user/dal/repo"
	"ims-server/internal/user/param"
	ioconsts "ims-server/pkg/consts"
	egoerror "ims-server/pkg/error"
	"ims-server/pkg/util"
)

type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, req *param.CreateUserRequest) (*param.CreateUserResponse, error) {
	user := &usermodel.User{
		Type:        req.Type,
		Account:     req.Account, // TODO: 写一个账号生成器
		Password:    req.Password,
		Name:        req.Name,
		Nickname:    req.Nickname,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Avatar:      req.Avatar,
		Status:      ioconsts.AccountStatusNormal,
	}
	err := userrepo.NewUserRepo().Create(ctx, user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}

	resp := userpack.ToUserResponse(user)
	return &param.CreateUserResponse{
		UserResponse: resp,
	}, nil
}

// GetUserByID 根据 ID 查询用户
func (s *userService) GetUserByID(ctx context.Context, req *param.GetUserByIDRequest) (*param.GetUserByIDResponse, error) {
	id := req.ID
	user, err := userrepo.NewUserRepo().Get(ctx, id)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	resp := userpack.ToUserResponse(user)
	return &param.GetUserByIDResponse{
		UserResponse: resp,
	}, nil
}

// MGetUserByID  根据用户 ID 列表获取多个用户信息
func (s *userService) MGetUserByID(ctx context.Context, req *param.MGetUserByIDRequest) (*param.MGetUserByIDResponse, error) {
	res, err := userrepo.NewUserRepo().MGet(ctx, req.IDs)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := []param.GetUserByIDResponse{}
	for _, user := range res {
		info := userpack.ToUserResponse(&user)
		resp = append(resp, param.GetUserByIDResponse{
			UserResponse: info,
		})
	}

	return &param.MGetUserByIDResponse{
		List: resp,
	}, nil
}

// UpdateUserById 根据 ID 更新用户
func (s *userService) UpdateUserByID(ctx context.Context, req *param.UpdateUserByIDRequest) (*param.UpdateUserByIDResponse, error) {
	_, err := userrepo.NewUserRepo().Get(ctx, req.ID) // 根据请求中的用户 ID 获取用户信息
	if err != nil {
		return nil, egoerror.ErrNotFound // 如果用户不存在，则返回错误信息
	}

	userMap := util.RequestToSnakeMapWithIgnoreZeroValueAndIDKey(req)

	update, err := userrepo.NewUserRepo().Update(ctx, req.ID, userMap) // 更新用户信息
	if err != nil {
		return nil, egoerror.ErrInvalidParam // 如果更新失败，则返回无效参数的错误信息
	}

	resp := userpack.ToUserResponse(update)
	return &param.UpdateUserByIDResponse{
		UserResponse: resp,
	}, nil
}

// DeleteUserByID 根据 ID 删除用户
func (s *userService) DeleteUserByID(ctx context.Context, r *param.DeleteUserByIDRequest) (*param.DeleteUserByIDResponse, error) {
	_, err := userrepo.NewUserRepo().Get(ctx, r.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	err = userrepo.NewUserRepo().Delete(ctx, r.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &param.DeleteUserByIDResponse{}, nil
}
