package userrepo

import (
	"context"
	"ims-server/internal/user/dal/model"
	ioconst "ims-server/pkg/consts"
	egoerror "ims-server/pkg/error"
	ioginx "ims-server/pkg/ginx"
)

type userRepo struct {
	ioginx.IRepo[usermodel.User]
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (r *userRepo) CreateEmptyAccount(ctx context.Context, password string, accountType ioconst.UserType) (*usermodel.User, error) {
	user := usermodel.User{
		Password: password,
		Type:     accountType,
	}
	err := r.Create(ctx, &user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	return &user, nil
}

func (r *userRepo) CreateUserAccount(ctx context.Context, name, password string, accountType ioconst.UserType) (*usermodel.User, error) {
	_, err := r.GetByUsername(ctx, name)
	if err == nil {
		return nil, egoerror.ErrRepeatedEntry
	}
	user := usermodel.User{
		Type:     accountType,
		Name:     name,
		Password: password,
		Email:    "",
	}
	err = r.Create(ctx, &user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*usermodel.User, error) {
	if username == "" {
		return nil, egoerror.ErrNotFound
	}
	var user usermodel.User
	err := r.DB().WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &user, nil
}
