package initialize

import (
	"encoding/gob"

	"github.com/go-webauthn/webauthn/webauthn"
)

func InitAll() {
	gob.Register(webauthn.SessionData{})

	initViper()
	initLogger()

	// initDB()
	initServer()
	initRouter()

	initAuth()

}
