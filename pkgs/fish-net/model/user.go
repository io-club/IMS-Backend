package model

import (
	"github.com/go-webauthn/webauthn/webauthn"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	AuthId      string
	AuthName    string
	DisplayName string
	Icon        string
	Credentials []webauthn.Credential `gorm:"-"`
}

func (user *User) WebAuthnID() []byte {
	return []byte(user.AuthId)
}

func (user *User) WebAuthnName() string {
	return user.AuthName
}

func (user *User) WebAuthnDisplayName() string {
	return user.DisplayName
}

func (user *User) WebAuthnIcon() string {
	return "https://pics.com/avatar.png"
}

func (user *User) WebAuthnCredentials() []webauthn.Credential {
	return user.Credentials
	// return []webauthn.Credential{}
}

func (user *User) AddWebAuthnCredential(credential webauthn.Credential) {
	user.Credentials = append(user.Credentials, credential)
}
