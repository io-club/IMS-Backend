package initialize

import (
	"fish_net/cmd/api/glb"

	"github.com/go-webauthn/webauthn/webauthn"
)

func initAuth() {
	Auth, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Go Webauthn",                                              // Display Name for your site
		RPID:          "localhost",                                                // Generally the FQDN for your site
		RPOrigins:     []string{"http://localhost:3000", "http://localhost:3001"}, // The origin URLs allowed for WebAuthn requests
		RPIcon:        "https://go-webauthn.local/logo.png",                       // Optional icon URL for your site
	})
	if err != nil {
		panic(err.Error())
	}
	glb.Auth = Auth
}
