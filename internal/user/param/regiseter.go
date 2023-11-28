package param

import ioconsts "ims-server/pkg/consts"

type SendVerification struct {
	Email string `json:"email" form:"email" binding:"required,email"`
	Url   string `json:"url" form:"url" binding:"required"`
}

type SendVerificationResponse struct {
}

type RegisterRequest struct {
	Type ioconsts.UserType `json:"type" form:"type" binding:"required"`

	Password string `json:"password" form:"password" binding:"required,gt=7"`

	Email            string `json:"email" form:"email" binding:"required,email"`
	VerificationCode string `json:"verificationCode" form:"verificationCode" binding:"required"`

	Nickname    string `json:"nickname" form:"nickname" binding:"lt=10"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`

	Avatar string `json:"avatar" form:"avatar"`
}

type RegisterResponse struct {
	UserResponse
}
