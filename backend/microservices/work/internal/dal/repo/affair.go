package repo

import (
	"ims-server/microservices/work/internal/dal/model"
	ioginx "ims-server/pkg/ginx"
)

var AffairSelect = []string{
	"id",
	"user_id",
	"topic",
	"content",
	"end_time",
	"created_at",
	"is_end",
}

type affairRepo struct {
	ioginx.IRepo[model.Affair]
}

func NewAffairRepo() *affairRepo {
	return &affairRepo{}
}
