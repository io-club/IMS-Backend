package usecase

import (
	"fishnet/domain"
	"fishnet/glb"
	"fishnet/user/repo"

	"go.uber.org/zap"
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
func (u *UserUsecase) CreateUser(users []*domain.User) (err error) {
	err = u.r.CreateUser(users)
	return
}

// Delete implements domain.UserUsecase
func (u *UserUsecase) DeleteUser(id int64) (err error) {
	err = u.r.DeleteUser(id)
	return
}

// Update implements domain.UserUsecase
func (u *UserUsecase) UpdateUser(id int64, nickName *string, icon *string) (err error) {
	err = u.r.UpdateUser(id, nickName, icon)
	return
}

// GetByID implements domain.UserUsecase
func (u *UserUsecase) QueryUser(userName *string, limit, offset int) ([]*domain.User, int64, error) {
	userModels, total, err := u.r.QueryUser(userName, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return userModels, total, nil
}

// GetByUsername implements domain.UserUsecase
func (u *UserUsecase) MGetUsers(userIDs []int64) ([]*domain.User, error) {
	glb.LOG.Info("field to query", zap.Int64s("userIDs", userIDs))
	users, err := u.r.MGetUsers(userIDs)
	if err != nil {
		return nil, err
	}
	return users, nil
}
