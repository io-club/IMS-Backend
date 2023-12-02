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
	user, err := repo.NewUserRepo().Get(ctx, id, repo.UserSelect...)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := pack.ToUserResponse(user)
	return &param.GetUserByIDResponse{
		UserResponse: resp,
	}, nil
}

// MGetUserByIDs  根据用户 ID 列表获取多个用户信息
func (u *userService) MGetUserByIDs(ctx context.Context, req *param.MGetUserByIDsRequest) (*param.MGetUserByIDsResponse, error) {
	res, err := repo.NewUserRepo().MGet(ctx, req.IDs, repo.UserSelect...)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := []param.UserResponse{}
	for _, user := range res {
		info := pack.ToUserResponse(&user)
		resp = append(resp, info)
	}

	return &param.MGetUserByIDsResponse{
		List: resp,
	}, nil
}

func (u *userService) GetUsers(ctx context.Context, req *param.GetUsersRequest) (*param.GetUsersResponse, error) {
	// 设置允许的过滤字段
	pageBuilder := req.Build(util.NewSet("id", "Name", "CreateAt"))

	total, err := repo.NewUserRepo().Count(ctx, pageBuilder)
	if err != nil {
		return nil, egoerror.ErrDBError
	}
	res, err := repo.NewUserRepo().Pageable(ctx, pageBuilder, repo.UserSelect...)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}

	resp := []param.UserResponse{}
	for _, user := range res {
		info := pack.ToUserResponse(&user)
		resp = append(resp, info)
	}

	return &param.GetUsersResponse{
		Total: total,
		List:  resp,
	}, nil
}

// UpdateUserByID 根据 ID 更新用户
func (u *userService) UpdateUserByID(ctx context.Context, req *param.UpdateUserByIDRequest) (*param.UpdateUserByIDResponse, error) {
	_, err := repo.NewUserRepo().Get(ctx, req.ID, repo.UserSelect...) // 根据请求中的用户 ID 获取用户信息
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
	_, err := repo.NewUserRepo().Get(ctx, r.ID, "id")
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	err = repo.NewUserRepo().Delete(ctx, r.ID)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &param.DeleteUserByIDResponse{}, nil
}
