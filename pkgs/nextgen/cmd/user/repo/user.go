package repo

import (
	"context"
	"fish_net/cmd/api/glb"
	"fish_net/cmd/user/domain"
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

func (u *userRepo) CreateUser(ctx context.Context, users []*domain.User) error {
	if err := glb.DB.Create(users).Error; err != nil {
		return err
	}
	return nil
}

func (u *userRepo) DeleteUser(ctx context.Context, userId int64) error {
	return glb.DB.Where("id = ?", userId).Delete(domain.User{}).Error
}

func (u *userRepo) UpdateUser(ctx context.Context, userID int64, nickName *string, avater *string) error {
	// params := map[string]interface{}{}
	params := &domain.User{}
	if nickName != nil {
		// params["nickname"] = *nickName
		params.Nickname = *nickName
	}
	if avater != nil {
		// params["avater"] = *icon
		params.Avater = *avater
	}
	return glb.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(params).Error
}

// QueryNote query list of note info
func (u *userRepo) QueryUser(ctx context.Context, userName *string, limit, offset int) ([]*domain.User, int64, error) {
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

func (u *userRepo) MGetUsers(ctx context.Context, userIDs []int64) ([]*domain.User, error) {

	var res []*domain.User
	if len(userIDs) == 0 {
		return res, nil
	}

	if err := glb.DB.Where("id in ?", userIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
