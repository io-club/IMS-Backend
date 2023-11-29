package service

import (
	"context"
	"ims-server/internal/user/bll/pack"
	"ims-server/internal/user/dal/repo"
	"ims-server/internal/user/param"
	egoerror "ims-server/pkg/error"
	"ims-server/pkg/util"
)

type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

// GetUserByID 根据 ID 查询用户
func (u *userService) GetUserByID(ctx context.Context, req *param.GetUserByIDRequest) (*param.GetUserByIDResponse, error) {
	id := req.ID
	user, err := repo.NewUserRepo().Get(ctx, id)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := pack.ToUserResponse(user)
	return &param.GetUserByIDResponse{
		UserResponse: resp,
	}, nil
}

// MGetUserByID  根据用户 ID 列表获取多个用户信息
func (u *userService) MGetUserByIDs(ctx context.Context, req *param.MGetUserByIDsRequest) (*param.MGetUserByIDsResponse, error) {
	res, err := repo.NewUserRepo().MGet(ctx, req.IDs)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := []param.GetUserByIDResponse{}
	for _, user := range res {
		info := pack.ToUserResponse(&user)
		resp = append(resp, param.GetUserByIDResponse{
			UserResponse: info,
		})
	}

	return &param.MGetUserByIDsResponse{
		List: resp,
	}, nil
}

// UpdateUserById 根据 ID 更新用户
func (u *userService) UpdateUserByID(ctx context.Context, req *param.UpdateUserByIDRequest) (*param.UpdateUserByIDResponse, error) {
	_, err := repo.NewUserRepo().Get(ctx, req.ID) // 根据请求中的用户 ID 获取用户信息
	if err != nil {
		return nil, egoerror.ErrNotFound // 如果用户不存在，则返回错误信息
	}

	userMap := util.RequestToSnakeMapWithIgnoreZeroValueAndIDKey(req)

	update, err := repo.NewUserRepo().Update(ctx, req.ID, userMap) // 更新用户信息
	if err != nil {
		return nil, egoerror.ErrInvalidParam // 如果更新失败，则返回无效参数的错误信息
	}

	resp := pack.ToUserResponse(update)
	return &param.UpdateUserByIDResponse{
		UserResponse: resp,
	}, nil
}

// DeleteUserByID 根据 ID 删除用户
func (u *userService) DeleteUserByID(ctx context.Context, r *param.DeleteUserByIDRequest) (*param.DeleteUserByIDResponse, error) {
	_, err := repo.NewUserRepo().Get(ctx, r.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	err = repo.NewUserRepo().Delete(ctx, r.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &param.DeleteUserByIDResponse{}, nil
}
