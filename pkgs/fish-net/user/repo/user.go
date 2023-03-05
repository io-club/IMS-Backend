package repo

import (
	"fishnet/domain"
	"fishnet/glb"

	"go.uber.org/zap"
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
	return glb.DB.Where("id = ?", userId).Delete(domain.User{}).Error
}

func (u *userRepo) UpdateUser(userID int64, nickName *string, icon *string) error {
	params := map[string]interface{}{}
	if nickName != nil {
		params["nickname"] = *nickName
	}
	if icon != nil {
		params["icon"] = *icon
	}
	return glb.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(params).Error
}

// QueryNote query list of note info
func (u *userRepo) QueryUser(userName *string, limit, offset int) ([]*domain.User, int64, error) {
	var total int64
	var res []*domain.User
	conn := glb.DB.Model(&domain.User{})

	if userName != nil {
		conn = conn.Where("username like ?", "%"+*userName+"%")
	}

	if err := conn.Count(&total).Error; err != nil {
		return res, total, err
	}

	if err := conn.Limit(limit).Offset(offset).Find(&res).Error; err != nil {
		return res, total, err
	}

	return res, total, nil
}

func (u *userRepo) MGetUsers(userIDs []int64) ([]*domain.User, error) {
	glb.LOG.Info("field to query", zap.Int64s("userIDs", userIDs))

	var res []*domain.User
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := glb.DB.Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	userNames := (func(users []*domain.User) (ret []string) {
		for _, user := range users {
			ret = append(ret, user.Username)
		}
		return
	})(res)
	glb.LOG.Info("query result", zap.Strings("userNames", userNames))
	return res, nil
}
