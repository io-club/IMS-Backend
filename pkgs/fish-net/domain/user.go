package domain

import (
	"encoding/binary"
	"fishnet/glb"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string
	Nickname    string
	Icon        string
	Credentials []webauthn.Credential `gorm:"-"`
}

// UserUsecase represent the User's usecases
type UserUsecase interface {
	Save(username string) (user *User, err error)
	Delete(id uint) (err error)
	Update(user *User) error
	GetByID(id uint) (user *User, err error)
	GetByUsername(username string) (user *User, err error)
}

// UserRepository represent the User's repository contract
type UserRepo interface {
	Save(username string) (*User, error)
	Delete(id uint) error
	Update(u *User) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
}

// go-webauthn implementation

// WebAuthnID returns the user's ID
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.Nickname
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return u.Icon
}

// AddCredential associates the credential to the user
func (u *User) AddWebAuthnCredential(cred webauthn.Credential) {
	glb.LOG.Info("[cred]: " + string(cred.ID))
	u.Credentials = append(u.Credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
