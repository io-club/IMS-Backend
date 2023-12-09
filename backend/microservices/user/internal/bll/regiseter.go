package bll

import (
	"context"
	"ims-server/microservices/user/internal/dal/model"
	"ims-server/microservices/user/internal/dal/repo"
	"ims-server/microservices/user/internal/job"
	"ims-server/microservices/user/internal/param"
	"ims-server/microservices/user/pkg/pack"
	ioconsts "ims-server/pkg/consts"
	egoerror "ims-server/pkg/error"
	iologger "ims-server/pkg/logger"
	ioredis "ims-server/pkg/redis"
	"ims-server/pkg/util"
	"time"
)

func (u *userService) SendVerification(ctx context.Context, req *param.SendVerification) (*param.SendVerificationResponse, error) {
	vcode, err := util.GetRandCode()
	if err != nil {
		return nil, egoerror.ErrFailedSend
	}
	// TODO：一定时间只能发送一份？
	// 先将验证码存入 redis
	rdb := ioredis.NewClient()
	err = rdb.Set(ctx, req.Email, vcode, 360*time.Second).Err()
	if err != nil {
		iologger.Error("redis set failed, err: %v", err)
		return nil, egoerror.ErrFailedSend
	}
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
	_, err := repo.NewUserRepo().GetByEmail(ctx, req.Email, "id")
	if err == nil {
		return nil, egoerror.ErrEmailExist
	}
	// 检查验证码是否正确
	code, err := ioredis.NewClient().Get(ctx, req.Email).Result()
	if err != nil {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	// 消费验证码
	err = ioredis.NewClient().Del(ctx, req.Email).Err()
	if err != nil {
		return nil, egoerror.ErrInvalidVerifyCode
	}

	if req.VerificationCode != code {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	// 加密密码
	req.Password, err = repo.NewUserRepo().EncryptedPassword(ctx, req.Password)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	user := &model.User{
		Type:        req.Type,
		Password:    req.Password,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		Status:      ioconsts.AccountStatusNormal,
	}
	err = repo.NewUserRepo().Create(ctx, user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	resp := pack.ToUserResponse(ctx, user)
	return &param.RegisterResponse{
		UserResponse: resp,
	}, nil
}

func (u *userService) NameLogin(ctx context.Context, req *param.NameLoginRequest) (*param.NameLoginResponse, error) {
	// 检查用户是否存在
	user, err := repo.NewUserRepo().GetByName(ctx, req.Name, repo.UserSelect...)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	// 检查密码是否正确
	user.Password, err = repo.NewUserRepo().DecryptedPassword(ctx, user.Password)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	if user.Password != req.Password {
		return nil, egoerror.ErrPasswordError
	}

	resp := pack.ToUserResponse(ctx, user)
	return &param.NameLoginResponse{
		UserResponse: resp,
	}, nil
}

func (u *userService) EmailLogin(ctx context.Context, req *param.EmailLoginRequest) (*param.EmailLoginResponse, error) {
	// 检查用户是否存在
	user, err := repo.NewUserRepo().GetByEmail(ctx, req.Email, repo.UserSelect...)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	// 检查验证码是否正确
	code, err := ioredis.NewClient().Get(ctx, req.Email).Result()
	if err != nil {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	// 消费验证码
	err = ioredis.NewClient().Del(ctx, req.Email).Err()
	if err != nil {
		return nil, egoerror.ErrInvalidVerifyCode
	}

	if req.VerificationCode != code {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	resp := pack.ToUserResponse(ctx, user)
	return &param.EmailLoginResponse{
		UserResponse: resp,
	}, nil
}

func (u *userService) RetrievePassword(ctx context.Context, req *param.RetrievePasswordRequest) (*param.RetrievePasswordResponse, error) {
	// 检查用户是否存在
	user, err := repo.NewUserRepo().GetByEmail(ctx, req.Email, repo.UserSelect...)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	// 检查验证码是否正确
	code, err := ioredis.NewClient().Get(ctx, req.Email).Result()
	if err != nil {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	if req.VerificationCode != code {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	// 加密密码
	req.Password, err = repo.NewUserRepo().EncryptedPassword(ctx, req.Password)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}

	m := map[string]interface{}{
		"password": req.Password,
	}
	_, err = repo.NewUserRepo().Update(ctx, user.ID, m)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}

	return &param.RetrievePasswordResponse{
		Name:     user.Name,
		Password: user.Password,
	}, nil
}
