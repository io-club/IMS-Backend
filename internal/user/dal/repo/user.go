package repo

import (
	"context"
	"ims-server/internal/user/dal/model"
	ioconst "ims-server/pkg/consts"
	egoerror "ims-server/pkg/error"
	ioginx "ims-server/pkg/ginx"
)

type userRepo struct {
	ioginx.IRepo[model.User]
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (r *userRepo) CreateEmptyAccount(ctx context.Context, password string, accountType ioconst.UserType) (*model.User, error) {
	user := model.User{
		Password: password,
		Type:     accountType,
	}
	err := r.Create(ctx, &user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	return &user, nil
}

func (r *userRepo) CreateUserAccount(ctx context.Context, name, password string, accountType ioconst.UserType) (*model.User, error) {
	_, err := r.GetByName(ctx, name)
	if err == nil {
		return nil, egoerror.ErrRepeatedEntry
	}
	user := model.User{
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

func (r *userRepo) GetByName(ctx context.Context, name string) (*model.User, error) {
	if name == "" {
		return nil, egoerror.ErrNotFound
	}
	var user model.User
	err := r.DB().WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &user, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.DB().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &user, nil
}
