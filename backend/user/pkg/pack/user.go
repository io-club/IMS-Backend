package pack

import (
	"ims-server/user/internal/dal/model"
	"ims-server/user/internal/param"
)

func ToUserResponse(u *model.User) param.UserResponse {
	res := param.UserResponse{
		ID:          u.ID,
		Type:        u.Type,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Avatar:      u.Avatar,
		Status:      u.Status,
	}
	return res
}
