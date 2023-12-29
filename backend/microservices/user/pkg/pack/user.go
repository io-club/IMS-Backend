package pack

import (
	"context"
	"fmt"
	"ims-server/microservices/user/internal/dal/model"
	"ims-server/microservices/user/internal/param"
	iooss "ims-server/pkg/oss"
	"time"
)

func ToUserResponse(ctx context.Context, u *model.User) param.UserResponse {
	res := param.UserResponse{
		ID:          u.ID,
		Type:        u.Type,
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Email:       u.Email,
		Status:      u.Status,
	}

	if u.Avatar != "" {
		if client, err := iooss.NewMinioClient(); err != nil {
			return res
		} else {
			signed, err := client.PresignedGetObject(ctx, iooss.DEFAULT_BUCKET_NAME, u.Avatar, 3*24*time.Hour, nil)
			if err != nil {
				return res
			}
			res.Avatar = fmt.Sprintf("%s", signed)
		}
	}
	return res
}
