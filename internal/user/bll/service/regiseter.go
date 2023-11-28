package service

import (
	"context"
	"ims-server/internal/user/bll/pack"
	"ims-server/internal/user/dal/model"
	"ims-server/internal/user/dal/repo"
	"ims-server/internal/user/job"
	"ims-server/internal/user/param"
	ioconsts "ims-server/pkg/consts"
	egoerror "ims-server/pkg/error"
	ioredis "ims-server/pkg/redis"
	"ims-server/pkg/util"
)

func (u *userService) SendVerification(ctx context.Context, req *param.SendVerification) (*param.SendVerificationResponse, error) {
	vcode, err := util.GetRandCode()
	if err != nil {
		return nil, egoerror.ErrFailedSend
	}
	// TODO：一定时间只能发送一份？
	// 先将验证码存入 redis
	rdb := ioredis.NewClient()
	rdb.Set(ctx, req.Email, vcode, 300)
	// 发送邮件
	err = job.SendEmail(req.Email, vcode, req.Url)
	if err != nil {
		return nil, egoerror.ErrFailedSend
	}
	return &param.SendVerificationResponse{}, nil
}

// Register 注册用户
func (u *userService) Register(ctx context.Context, req *param.RegisterRequest) (*param.RegisterResponse, error) {
	// 检查该邮箱是否已被注册
	_, err := repo.NewUserRepo().GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, egoerror.ErrEmailExist
	}
	// 检查验证码是否正确
	code := ioredis.NewClient().Get(ctx, req.Email)
	if req.VerificationCode != code.String() {
		return nil, egoerror.ErrInvalidVerifyCode
	}

	user := &model.User{
		Type:        req.Type,
		Account:     req.Email,
		Password:    req.Password,
		Name:        req.Email,
		Nickname:    req.Nickname,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Avatar:      req.Avatar,
		Status:      ioconsts.AccountStatusNormal,
	}
	err = repo.NewUserRepo().Create(ctx, user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	resp := pack.ToUserResponse(user)
	return &param.RegisterResponse{
		UserResponse: resp,
	}, nil
}
