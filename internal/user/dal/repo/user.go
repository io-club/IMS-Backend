package repo

import (
	"context"
	"ims-server/internal/user/dal/model"
	ioconfig "ims-server/pkg/config"
	ioconst "ims-server/pkg/consts"
	"ims-server/pkg/encryption"
	egoerror "ims-server/pkg/error"
	ioginx "ims-server/pkg/ginx"
	"log"
)

var userSelect = []string{
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

func (r *userRepo) EncryptedPassword(ctx context.Context, password string) (string, error) {
	encrypt, err := encryption.Encrypt([]byte(password), encryption.AES, ioconfig.GetEncryptionConf().AesKey)
	if err != nil {
		log.Printf("encrypt failed, err: %v", err)
		return "", err
	}
	return string(encrypt), nil
}

func (r *userRepo) DecryptedPassword(ctx context.Context, password string) (string, error) {
	decrypt, err := encryption.Decrypt([]byte(password), encryption.AES, ioconfig.GetEncryptionConf().AesKey)
	if err != nil {
		log.Printf("encrypt failed, err: %v", err)
		return "", err
	}
	return string(decrypt), nil
}

func (r *userRepo) CreateEmptyAccount(ctx context.Context, password string, accountType ioconst.UserType) (*model.User, error) {
	user := model.User{
		Password: password,
		Type:     accountType,
	}
	err := r.Create(ctx, &user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	return &user, nil
}

func (r *userRepo) CreateUserAccount(ctx context.Context, name, password string, accountType ioconst.UserType) (*model.User, error) {
	_, err := r.GetByName(ctx, name)
	if err == nil {
		return nil, egoerror.ErrRepeatedEntry
	}
	user := model.User{
		Type:     accountType,
		Name:     name,
		Password: password,
		Email:    "",
	}
	err = r.Create(ctx, &user)
	if err != nil {
		return nil, egoerror.ErrInvalidParam
	}
	return &user, nil
}

func (r *userRepo) GetByName(ctx context.Context, name string) (*model.User, error) {
	if name == "" {
		return nil, egoerror.ErrNotFound
	}
	var user model.User
	err := r.DB().WithContext(ctx).Select(append(userSelect, "password")).Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &user, nil
}

func (u *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := u.DB().WithContext(ctx).Select(append(userSelect, "password")).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, egoerror.ErrNotFound
	}
	return &user, nil
}
