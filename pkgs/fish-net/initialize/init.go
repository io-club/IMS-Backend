package initialize

import (
	"encoding/gob"

	"github.com/go-webauthn/webauthn/webauthn"
)

func InitAll() {
	gob.Register(webauthn.SessionData{})

	// load config
	initViper()
	// init logger
	initLogger()

	// db
	initDB()

	// passkey webauthn
	initAuth()

}
