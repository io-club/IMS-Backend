package domain

import (
	"encoding/binary"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type WebAuthnCredential struct {
	gorm.Model
	UserID          int64 `gorm:"index"`
	CredentialID    []byte
	PublicKey       []byte
	AttestationType string
	AAGUID          []byte
	SignCount       uint32
	CloneWarning    bool
}

func (cred WebAuthnCredential) TableName() string {
	return "webauthn_credential"
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

type WebAuthnCredentialUsecase interface {
	QueryCredential(userID int64) []webauthn.Credential
	QueryByPublicKey(publicKey []byte) (*WebAuthnCredential, error)
	QueryByCredentialID(credentialID []byte) (*WebAuthnCredential, error)
	CreateCredential(userID int64, cred *webauthn.Credential) (*WebAuthnCredential, error)
}
