package pack

import (
	"ims-server/internal/user/dal/model"
	"ims-server/internal/user/param"
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
