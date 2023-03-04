package repo

import (
	"fishnet/domain"
	"fishnet/glb"

	"gorm.io/gorm"
)

type userRepo struct {
}

var _userRepo *userRepo

func NewUserRepo() domain.UserRepo {
	if _userRepo == nil {
		_userRepo = &userRepo{}
	}
	return _userRepo
}

func (u *userRepo) Save(username string) (user *domain.User, err error) {
	defaultIcon := "https://pics.com/avatar.png"
	user = &domain.User{
		Username: username,
		Nickname: username,
		Icon:     defaultIcon,
		// Credentials: []webauthn.Credential{},
	}
	err = glb.DB.Save(user).Error
	return
}

func (u *userRepo) Delete(id uint) (err error) {
	err = glb.DB.Delete(&domain.User{}, id).Error
	return
}

func (u *userRepo) Update(user *domain.User) (err error) {
	err = glb.DB.Save(user).Error
	return
}

func (u *userRepo) GetByID(id uint) (user *domain.User, err error) {
	user = &domain.User{Model: gorm.Model{ID: id}}
	err = glb.DB.First(&user).Error
	return
}

func (u *userRepo) GetByUsername(username string) (user *domain.User, err error) {
	err = glb.DB.Where(&domain.User{Username: username}).First(&user).Error
	return
}
