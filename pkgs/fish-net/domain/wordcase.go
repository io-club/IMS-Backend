package domain

import "gorm.io/gorm"

type Wordcase struct {
	gorm.Model
	GroupName string // 分组名称 例如：device, express
	Key       string // 键 例如：sensor, expressNo
	Value     string // 值 例如：temperature, {name: "JD", code: "JD"}
	Order     int    // 顺序
	Disable   bool   // 是否禁用
	Remark    string // 备注
}

func (Wordcase) TableName() string {
	return "wordcase"
}

type WordcaseRepo interface {
	CreateWordcase(wordcases []*Wordcase) error
	DeleteWordcase(wordcaseID int64) error
	UpdateWordcase(wordcaseID int64, value *string, order *int, disable *bool, remark *string) error
	QueryWordcase(wordcaseID *int64, groupName, key *string, limit, offset int) ([]*Wordcase, error)
	MGetWordcases(wordcaseIDs []int64) ([]*Wordcase, error)
}

type WordcaseUsecase interface {
	CreateWordcase(wordcase []*Wordcase) error
	DeleteWordcase(wordcaseID int64) error
	UpdateWordcase(wordcaseID int64, value *string, order *int, disable *bool, remark *string) error
	QueryWordcase(wordcaseID *int64, groupName, key *string, limit, offset int) ([]*Wordcase, error)
	MGetWordcases(wordcaseIDs []int64) ([]*Wordcase, error)
}
