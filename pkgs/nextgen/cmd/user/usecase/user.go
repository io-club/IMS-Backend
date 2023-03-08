package usecase

import (
	"context"
	"fish_net/cmd/user/domain"
	"fish_net/cmd/user/repo"
)

type UserUsecase struct {
	r domain.UserRepo
}

var _userUsecase *UserUsecase

func NewUserUsecase() domain.UserUsecase {
	if _userUsecase == nil {
		_userUsecase = &UserUsecase{
			r: repo.NewUserRepo(),
		}
	}
	return _userUsecase
}

// Save implements domain.UserUsecase
func (u *UserUsecase) CreateUser(ctx context.Context, users []*domain.User) (err error) {
	err = u.r.CreateUser(ctx, users)
	return
}

// Delete implements domain.UserUsecase
func (u *UserUsecase) DeleteUser(ctx context.Context, id int64) (err error) {
	err = u.r.DeleteUser(ctx, id)
	return
}

// Update implements domain.UserUsecase
func (u *UserUsecase) UpdateUser(ctx context.Context, id int64, nickName *string, icon *string) (err error) {
	err = u.r.UpdateUser(ctx, id, nickName, icon)
	return
}

// GetByID implements domain.UserUsecase
func (u *UserUsecase) QueryUser(ctx context.Context, userName *string, limit, offset int) ([]*domain.User, int64, error) {
	userModels, total, err := u.r.QueryUser(ctx, userName, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return userModels, total, nil
}

// GetByUsername implements domain.UserUsecase
func (u *UserUsecase) MGetUsers(ctx context.Context, userIDs []int64) ([]*domain.User, error) {
	users, err := u.r.MGetUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return users, nil
}
