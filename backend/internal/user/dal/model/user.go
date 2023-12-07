package model

import (
	"gorm.io/gorm"
	ioconsts "ims-server/pkg/consts"
)

// User 用户
type User struct {
	gorm.Model
	Type ioconsts.UserType `gorm:"type:char;size:15;comment:账户类型;not null"`

	Name     string `gorm:"type:char;size:10;comment:用户名;unique"`
	Password string `gorm:"type:string;comment:密码;not null"`
	Email    string `gorm:"type:char;size:30;comment:邮箱;unique"`

	PhoneNumber string `gorm:"type:char;size:20;comment:手机号"`
	Avatar      string `gorm:"type:string;comment:头像"`

	Status ioconsts.AccountStatus `gorm:"type:char;size:10;comment:用户状态"`
}

func (User) TableName() string {
	return "ims_user_user"
}
