package glb

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	// MG
	Auth *webauthn.WebAuthn
	// VP   *viper.Viper
)
