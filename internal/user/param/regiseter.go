package param

import ioconsts "ims-server/pkg/consts"

type LoginResponse struct {
	UserResponse
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SendVerification struct {
	Email string `json:"email" form:"email" binding:"required,email"`
	Url   string `json:"url" form:"url" binding:"required"`
}

type SendVerificationResponse struct {
}

type RegisterRequest struct {
	Type ioconsts.UserType `json:"type" form:"type" binding:"required"`

	Name     string `json:"name" form:"name" binding:"required,lt=10"`
	Password string `json:"password" form:"password" binding:"required,gt=7,password"`

	Email            string `json:"email" form:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode" form:"verificationCode" binding:"required"`

	PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`

	Avatar string `json:"avatar" form:"avatar"`
}

type RegisterResponse struct {
	UserResponse
}

type NameLoginRequest struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type NameLoginResponse struct {
	LoginResponse
}

type EmailLoginRequest struct {
	Email            string `json:"email" form:"email" binding:"required"`
	VerificationCode string `json:"verificationCode" form:"verificationCode" binding:"required"`
}

type EmailLoginResponse struct {
	LoginResponse
}

type RetrievePasswordRequest struct {
	Email            string `json:"email" form:"email" binding:"required"`
	VerificationCode string `json:"verificationCode" form:"verificationCode" binding:"required"`
	Password         string `json:"password" form:"password" binding:"required,gt=7,password"`
}

type RetrievePasswordResponse struct {
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}
