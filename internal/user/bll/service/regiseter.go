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
	rdb.Set(ctx, req.Email, vcode, 360*time.Second)
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
	code, err := ioredis.NewClient().Get(ctx, req.Email).Result()
	if err != nil {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	if req.VerificationCode != code {
		return nil, egoerror.ErrInvalidVerifyCode
	}
	// TODO: 加密密码
	user := &model.User{
		Type:        req.Type,
		Password:    req.Password,
		Name:        req.Name,
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

func (u *userService) NameLogin(ctx context.Context, req *param.NameLoginRequest) (*param.NameLoginResponse, error) {
	// 检查用户是否存在
	user, err := repo.NewUserRepo().GetByName(ctx, req.Name)
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	// 检查密码是否正确
	if user.Password != req.Password {
		return nil, egoerror.ErrPasswordError
	}
	resp := pack.ToUserResponse(user)

	payload := util.JwtPayload{
		Issue:       "ims",
		Audience:    "ims",
		Subject:     "refresh-Token",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(),
		UserDefined: map[string]any{"uid": user.ID, "uname": user.Name, "utype": user.Type},
	}
	refreshToken, err := util.GenJwt(util.DefautHeader, payload)
	// 生成 access token
	payload.Subject = "access-Token"
	payload.Expiration = time.Now().Add(1 * time.Hour).Unix()
	accessToken, err := util.GenJwt(util.DefautHeader, payload)
	return &param.NameLoginResponse{
		LoginResponse: param.LoginResponse{
			UserResponse: resp,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (u *userService) EmailLogin(ctx context.Context, req *param.EmailLoginRequest) (*param.EmailLoginResponse, error) {
	// 检查用户是否存在
	user, err := repo.NewUserRepo().GetByEmail(ctx, req.Email)
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
	resp := pack.ToUserResponse(user)

	payload := util.JwtPayload{
		Issue:       "ims",
		Audience:    "ims",
		Subject:     "refresh-Token",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(3 * 24 * time.Hour).Unix(),
		UserDefined: map[string]any{"uid": user.ID, "uname": user.Name, "utype": user.Type},
	}
	refreshToken, err := util.GenJwt(util.DefautHeader, payload)
	// 生成 access token
	payload.Subject = "access-Token"
	payload.Expiration = time.Now().Add(1 * time.Hour).Unix()
	accessToken, err := util.GenJwt(util.DefautHeader, payload)
	return &param.EmailLoginResponse{
		LoginResponse: param.LoginResponse{
			UserResponse: resp,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (u *userService) RetrievePassword(ctx context.Context, req *param.RetrievePasswordRequest) (*param.RetrievePasswordResponse, error) {
	// 检查用户是否存在
	user, err := repo.NewUserRepo().GetByEmail(ctx, req.Email)
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
