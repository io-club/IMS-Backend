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

	if client, err := iooss.NewMinioClient(); err != nil {
		res.Avatar = ""
		return res
	} else {
		//reqParams := make(url.Values)
		//reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.jpg\"")
		signed, err := client.PresignedGetObject(ctx, iooss.DefaultBucketName, u.Avatar, 3*24*time.Hour, nil)
		if err != nil {
			res.Avatar = ""
		}
		res.Avatar = fmt.Sprintf("%s", signed)
	}
	return res
}
