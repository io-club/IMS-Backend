package repo

import (
	"context"
	"ims-server/microservices/user/internal/dal/model"
	ioconfig "ims-server/pkg/config"
	ioconst "ims-server/pkg/consts"
	"ims-server/pkg/encryption"
	ioginx "ims-server/pkg/ginx"
	iologger "ims-server/pkg/logger"
)

var UserSelect = []string{
	"id",
	"type",
	"name",
	"email",
	"phone_number",
	"avatar",
	"status",
}

type userRepo struct {
	ioginx.IRepo[model.User]
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (u *userRepo) EncryptedPassword(ctx context.Context, password string) (string, error) {
	encrypt, err := encryption.Encrypt([]byte(password), encryption.AES, ioconfig.GetEncryptionConf().AesKey)
	if err != nil {
		iologger.Warn("encrypt failed, err: %v", err)
		return "", err
	}
	return string(encrypt), nil
}

func (u *userRepo) DecryptedPassword(ctx context.Context, password string) (string, error) {
	decrypt, err := encryption.Decrypt([]byte(password), encryption.AES, ioconfig.GetEncryptionConf().AesKey)
	if err != nil {
		iologger.Warn("encrypt failed, err: %v", err)
		return "", err
	}
	return string(decrypt), nil
}

func (u *userRepo) CreateEmptyAccount(ctx context.Context, password string, accountType ioconst.UserType) (*model.User, error) {
	user := model.User{
		Password: password,
		Type:     accountType,
	}
	err := u.Create(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) CreateUserAccount(ctx context.Context, name, password string, accountType ioconst.UserType) (*model.User, error) {
	_, err := u.GetByName(ctx, name)
	if err == nil {
		return nil, err
	}
	user := model.User{
		Type:     accountType,
		Name:     name,
		Password: password,
		Email:    "",
	}
	err = u.Create(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) GetByName(ctx context.Context, name string, selectField ...string) (*model.User, error) {
	var user model.User
	selectField = append(selectField, "password")
	err := u.DB().WithContext(ctx).Select(append(UserSelect, selectField...)).Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string, selectField ...string) (*model.User, error) {
	var user model.User
	selectField = append(selectField, "password")
	err := u.DB().WithContext(ctx).Select(append(UserSelect, selectField...)).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
