package repo

import (
	"fishnet/common/consts"
	"fishnet/domain"
	"fishnet/glb"
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

func (u *userRepo) CreateUser(users []*domain.User) error {
	if err := glb.DB.Create(users).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) DeleteUser(userId int64) error {
	return glb.DB.Where("id = ?", userId).Delete(&domain.User{}).Error
}

func (u *userRepo) UpdateUser(userID int64, nickName *string, icon *string) error {
	params := map[string]interface{}{}
	if nickName != nil && *nickName != "" {
		params["nickname"] = *nickName
	}
	if icon != nil && *icon != "" {
		params["icon"] = *icon
	}
	return glb.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(params).Error
}

// QueryNote query list of note info
func (u *userRepo) QueryUser(userID *int64, userName *string, nickName *string, limit, offset int) ([]*domain.User, error) {
	var total int64
	var res []*domain.User
	conn := glb.DB.Model(&domain.User{})
	if userID != nil && *userID > 0 {
		conn = conn.Where("id = ?", *userID)
	} else {
		if userName != nil && *userName != "" {
			conn = conn.Where("username = ?", *userName)
		}
		if nickName != nil && *nickName != "" {
			conn = conn.Where("nickname = ?", *nickName)
		}
		if limit == 0 {
			limit = consts.DefaultLimit
		}
		conn = conn.Limit(limit).Offset(offset)
		if err := conn.Count(&total).Error; err != nil {
			return res, err
		}
	}
	if err := conn.Find(&res).Order("id desc").Error; err != nil {
		return res, err
	}

	return res, nil
}

func (u *userRepo) MGetUsers(userIDs []int64) ([]*domain.User, error) {
	var res []*domain.User
	if len(userIDs) == 0 {
		return res, nil
	}
	if err := glb.DB.Where("id in (?)", userIDs).Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}
