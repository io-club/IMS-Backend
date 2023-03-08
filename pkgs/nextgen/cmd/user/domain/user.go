package domain

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"index"`
	Nickname    string
	Avater      string
	Credentials []webauthn.Credential `gorm:"-"`
}

func (User) TableName() string {
	return "user"
}

// UserRepository represent the User's repository contract
type UserRepo interface {
	CreateUser(ctx context.Context, users []*User) error
	DeleteUser(ctx context.Context, userID int64) error
	UpdateUser(ctx context.Context, userID int64, nickName *string, icon *string) error
	QueryUser(ctx context.Context, userName *string, limit, offset int) ([]*User, int64, error)
	MGetUsers(ctx context.Context, userIDs []int64) ([]*User, error)
}

// UserUsecase represent the User's usecases
type UserUsecase interface {
	CreateUser(ctx context.Context, users []*User) error
	DeleteUser(ctx context.Context, userID int64) error
	UpdateUser(ctx context.Context, userID int64, nickName *string, icon *string) error
	QueryUser(ctx context.Context, userName *string, limit, offset int) ([]*User, int64, error)
	MGetUsers(ctx context.Context, userIDs []int64) ([]*User, error)
}
