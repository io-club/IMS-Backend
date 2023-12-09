package bll

import (
	"context"
	"ims-server/microservices/user/internal/dal/repo"
	"ims-server/microservices/user/internal/param"
	"ims-server/microservices/user/pkg/pack"
	egoerror "ims-server/pkg/error"
	iologger "ims-server/pkg/logger"
	iooss "ims-server/pkg/oss"
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

	resp := pack.ToUserResponse(ctx, user)
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
		info := pack.ToUserResponse(ctx, &user)
		resp = append(resp, info)
	}

	return &param.MGetUserByIDsResponse{
		List: resp,
	}, nil
}

func (u *userService) GetUsers(ctx context.Context, req *param.GetUsersRequest) (*param.GetUsersResponse, error) {
	// 设置允许的过滤字段
	pageBuilder := req.Build("id", "Name", "CreateAt")

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
		info := pack.ToUserResponse(ctx, &user)
		resp = append(resp, info)
	}

	return &param.GetUsersResponse{
		Total: total,
		List:  resp,
	}, nil
}

// UpdateUserByID 根据 ID 更新用户
func (u *userService) UpdateUserByID(ctx context.Context, req *param.UpdateUserByIDRequest) (*param.UpdateUserByIDResponse, error) {
	// 确保用户存在
	_, err := repo.NewUserRepo().Get(ctx, req.ID, repo.UserSelect...) // 根据请求中的用户 ID 获取用户信息
	if err != nil {
		return nil, egoerror.ErrNotFound // 如果用户不存在，则返回错误信息
	}
	// 只有用户本身才能修改自己的头像
	userID := ctx.Value("uid").(float64)
	if req.ID != uint(userID) {
		return nil, egoerror.ErrNotPermitted
	}

	userMap := util.RequestToSnakeMapWithIgnoreZeroValueAndIDKey(req)

	update, err := repo.NewUserRepo().Update(ctx, req.ID, userMap) // 更新用户信息
	if err != nil {
		return nil, egoerror.ErrInvalidParam // 如果更新失败，则返回无效参数的错误信息
	}

	resp := pack.ToUserResponse(ctx, update)
	return &param.UpdateUserByIDResponse{
		UserResponse: resp,
	}, nil
}

func (u *userService) UploadAvatar(ctx context.Context, req *param.UploadAvatarRequest) (*param.UploadAvatarResponse, error) {
	// 只有用户本身才能修改自己的头像
	userID := ctx.Value("uid").(float64)
	if req.ID != uint(userID) {
		return nil, egoerror.ErrNotPermitted
	}
	user, err := repo.NewUserRepo().Get(ctx, req.ID, "avatar")
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	// Upload avatar
	client, err := iooss.NewMinioClient()
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	fileName, err := client.PutObject(ctx, iooss.DefaultBucketName, req.Avatar.Filename, req.Avatar)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	// Update user avatar
	_, err = repo.NewUserRepo().Update(ctx, req.ID, map[string]interface{}{
		"avatar": fileName,
	})
	if err != nil {
		return nil, egoerror.ErrInvalidParam // 如果更新失败，则返回无效参数的错误信息
	}
	// 更新成功则删除老图片
	if user.Avatar != "" {
		err = client.DeleteObject(ctx, iooss.DefaultBucketName, user.Avatar)
		if err != nil {
			iologger.Info("delete old avatar failed, fileName: %s, err: %s", user.Avatar, err.Error())
		}
	}
	return &param.UploadAvatarResponse{}, nil
}

// DeleteUserByID 根据 ID 删除用户
func (u *userService) DeleteUserByID(ctx context.Context, r *param.DeleteUserByIDRequest) (*param.DeleteUserByIDResponse, error) {
	// TODO：删除用户后同时删除头像？
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
