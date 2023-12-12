package pack

import (
	"ims-server/microservices/work/internal/dal/model"
	"ims-server/microservices/work/internal/param"
)

func ToAffairResponse(affair model.Affair) param.AffairResponse {
	return param.AffairResponse{
		ID:        affair.ID,
		Topic:     affair.Topic,
		Context:   affair.Content,
		StartTime: affair.CreatedAt,
		EndTime:   affair.EndTime,
		IsEnd:     affair.IsEnd,
	}
}
