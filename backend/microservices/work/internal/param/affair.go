package param

import "time"

type AffairResponse struct {
	ID        uint      `json:"id" form:"id"`
	Topic     string    `json:"topic" form:"topic"`
	Context   string    `json:"context" form:"context"`
	StartTime time.Time `json:"startTime" form:"startTime"`
	EndTime   time.Time `json:"endTime" form:"endTime"`
	IsEnd     bool      `json:"isEnd" form:"isEnd"`
}
