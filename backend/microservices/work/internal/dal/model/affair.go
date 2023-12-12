package model

import (
	"gorm.io/gorm"
	"time"
)

type Affair struct {
	gorm.Model
	UserID uint `gorm:"type:uint;comment:用户 ID;not null"`

	Topic   string `gorm:"type:char(20);comment:标题;not null"`
	Content string `gorm:"type:text;comment:事件内容"`

	EndTime time.Time `gorm:"type:datetime;comment:结束时间"`

	IsEnd bool `gorm:"type:bool;comment:是否结束"`
}

func (Affair) TableName() string {
	return "ims_work_affair"
}
