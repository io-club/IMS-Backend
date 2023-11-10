package http

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/user/usecase"
)

const REGISTER_SESSION_DATA_KEY = "register_session_data"
const LOGIN_SESSION_DATA_KEY = "login_session_data"

var _userUsecase domain.UserUsecase
var _webAuthnCredentialUsecase domain.WebAuthnCredentialUsecase

func init() {
	_userUsecase = usecase.NewUserUsecase()
	_webAuthnCredentialUsecase = usecase.NewWebAuthnCredentialUsecase()
}
