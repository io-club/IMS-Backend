package usecase

import (
	"fishnet/domain"
	"fishnet/user/repo"
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

// Delete implements domain.UserUsecase
func (u *UserUsecase) Delete(id uint) (err error) {
	err = u.r.Delete(id)
	return
}

// GetByID implements domain.UserUsecase
func (u *UserUsecase) GetByID(id uint) (user *domain.User, err error) {
	user, err = u.r.GetByID(id)
	return
}

// GetByUsername implements domain.UserUsecase
func (u *UserUsecase) GetByUsername(username string) (user *domain.User, err error) {
	user, err = u.r.GetByUsername(username)
	return
}

// Save implements domain.UserUsecase
func (u *UserUsecase) Save(username string) (user *domain.User, err error) {
	user, err = u.r.Save(username)
	return
}

// Update implements domain.UserUsecase
func (u *UserUsecase) Update(user *domain.User) (err error) {
	err = u.r.Update(user)
	return
}
