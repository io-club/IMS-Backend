package initialize

import (
	"encoding/gob"

	"github.com/go-webauthn/webauthn/webauthn"
)

func init() {
	gob.Register(webauthn.SessionData{})

	// load config
	// initViper()
	// init logger
	// initLogger()

	// db
	initDB()
	initLogger()

	// server
	// initServer()
	// router
	// initRouter()

	// passkey webauthn
	// initAuth()

}
