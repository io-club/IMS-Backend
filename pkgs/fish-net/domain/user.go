package domain

import (
	"encoding/binary"
	"fishnet/glb"

	"github.com/go-webauthn/webauthn/protocol"
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

// UserRepository represent the User's repository contract
type UserRepo interface {
	CreateUser(users []*User) error
	DeleteUser(userID int64) error
	UpdateUser(userID int64, nickName *string, icon *string) error
	QueryUser(userName *string, limit, offset int) ([]*User, int64, error)
	MGetUsers(userIDs []int64) ([]*User, error)
}

// UserUsecase represent the User's usecases
type UserUsecase interface {
	CreateUser(users []*User) error
	DeleteUser(userID int64) error
	UpdateUser(userID int64, nickName *string, icon *string) error
	QueryUser(userName *string, limit, offset int) ([]*User, int64, error)
	MGetUsers(userIDs []int64) ([]*User, error)
}

type Credential struct {
	gorm.Model
	UserID       int64
	CredentialID int64
	PublicKey    string
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

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all a user's credentials
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.Credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
