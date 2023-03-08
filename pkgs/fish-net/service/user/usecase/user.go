package usecase

import (
	"errors"
	"fishnet/domain"
	"fishnet/glb"
	"fishnet/service/user/repo"

	"go.uber.org/zap"
)

type userUsecase struct {
	r domain.UserRepo
}

var _userUsecase *userUsecase

func NewUserUsecase() domain.UserUsecase {
	if _userUsecase == nil {
		_userUsecase = &userUsecase{
			r: repo.NewUserRepo(),
		}
	}
	return _userUsecase
}

// Save implements domain.UserUsecase
func (u *userUsecase) CreateUser(users []*domain.User) (err error) {
	err = u.r.CreateUser(users)
	return
}

// Delete implements domain.UserUsecase
func (u *userUsecase) DeleteUser(id int64) (err error) {
	err = u.r.DeleteUser(id)
	return
}

// Update implements domain.UserUsecase
func (u *userUsecase) UpdateUser(id int64, nickName *string, icon *string) (err error) {
	err = u.r.UpdateUser(id, nickName, icon)
	return
}

// GetByID implements domain.UserUsecase
func (u *userUsecase) QueryUser(userID *int64, userName *string, nickName *string, limit, offset int) ([]*domain.User, error) {
	userModels, err := u.r.QueryUser(userID, userName, nickName, limit, offset)
	if err != nil {
		return nil, err
	}
	return userModels, nil
}

// GetByUsername implements domain.UserUsecase
func (u *userUsecase) MGetUsers(userIDs []int64) ([]*domain.User, error) {
	glb.LOG.Info("field to query", zap.Int64s("userIDs", userIDs))
	users, err := u.r.MGetUsers(userIDs)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CheckUserExist implements domain.UserUsecase
func (u *userUsecase) CheckUserExist(userID *int64, userName *string) (bool, error) {
	users, err := u.QueryUser(userID, userName, nil, 1, 0)
	if err != nil || len(users) == 0 {
		return false, err
	}
	return true, nil
}

// FindByID implements domain.UserUsecase
func (u *userUsecase) FindByID(userID *int64) (*domain.User, error) {
	users, err := u.QueryUser(userID, nil, nil, 1, 0)
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}
