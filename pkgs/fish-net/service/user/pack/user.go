package pack

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/common"
)

type CreateUserRequest struct {
	Username string `json:"username"`
}

type GetUserResponse struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Icon     string `json:"icon"`
}

func User(u *domain.User) GetUserResponse {
	return GetUserResponse{
		Username: u.Username,
		Nickname: u.Nickname,
		Icon:     u.Icon,
	}
}

func Users(us []*domain.User) []GetUserResponse {
	var ret []GetUserResponse
	for _, u := range us {
		ret = append(ret, User(u))
	}
	return ret
}

type QueryUserRequest struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	common.PageRequest
}

type UpdateUserRequest struct {
	Nickname *string `json:"nickname"`
	Icon     *string `json:"icon"`
}
