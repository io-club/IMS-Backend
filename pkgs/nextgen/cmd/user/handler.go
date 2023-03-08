package main

import (
	"context"
	"fish_net/cmd/user/domain"
	"fish_net/cmd/user/pack"
	"fish_net/cmd/user/usecase"
	user "fish_net/kitex_gen/user"
	"fish_net/pkg/consts"
	"fish_net/pkg/errno"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp = new(user.CreateUserResponse)

	if err = req.IsValid(); err != nil {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = usecase.NewUserUsecase().CreateUser(ctx, []*domain.User{
		{
			Username: req.Username,
			Nickname: req.Username,
			Avater:   "https://pics.com/avatar.png",
		},
	})
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// DeleteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (resp *user.DeleteUserResponse, err error) {
	resp = new(user.DeleteUserResponse)

	if err = req.IsValid(); err != nil {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = usecase.NewUserUsecase().DeleteUser(ctx,
		req.UserId,
	)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// UpdateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	resp = new(user.UpdateUserResponse)

	if err = req.IsValid(); err != nil {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	err = usecase.NewUserUsecase().UpdateUser(ctx,
		req.UserId,
		req.Nickname,
		req.Avater,
	)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	return resp, nil
}

// QueryUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) QueryUser(ctx context.Context, req *user.QueryUserRequest) (resp *user.QueryUserResponse, err error) {
	resp = new(user.QueryUserResponse)

	if err = req.IsValid(); err != nil {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	if req.Limit == 0 {
		req.Limit = consts.DefaultLimit
	}

	users, total, err := usecase.NewUserUsecase().QueryUser(ctx,
		req.Username,
		int(req.Limit),
		int(req.Offset),
	)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.Users = pack.Users(users)
	resp.Total = int32(total)
	return resp, nil
}

// MGetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) MGetUser(ctx context.Context, req *user.MGetUserRequest) (resp *user.MGetUserResponse, err error) {
	resp = new(user.MGetUserResponse)

	if err = req.IsValid(); err != nil {
		resp.BaseResp = pack.BuildBaseResp(errno.ParamErr)
		return resp, nil
	}

	users, err := usecase.NewUserUsecase().MGetUsers(ctx, req.UserIds)
	if err != nil {
		resp.BaseResp = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.BaseResp = pack.BuildBaseResp(errno.Success)
	resp.Users = pack.Users(users)
	return resp, nil
}
