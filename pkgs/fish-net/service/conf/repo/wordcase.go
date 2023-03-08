package repo

import (
	"fishnet/domain"
	"fishnet/glb"
)

var _wordcaseRepo domain.WordcaseRepo

type wordcaseRepo struct {
}

func NewWordcaseRepo() domain.WordcaseRepo {
	if _wordcaseRepo == nil {
		_wordcaseRepo = &wordcaseRepo{}
	}
	return _wordcaseRepo
}

func (u *wordcaseRepo) CreateWordcase(wordcases []*domain.Wordcase) error {
	if err := glb.DB.Create(wordcases).Error; err != nil {
		return err
	}
	return nil
}

func (u *wordcaseRepo) DeleteWordcase(wordcaseId int64) error {
	return glb.DB.Where("id = ?", wordcaseId).Delete(domain.Wordcase{}).Error
}

// UpdateWordcase implements domain.WordcaseRepo
func (*wordcaseRepo) UpdateWordcase(wordcaseID int64, value *string, order *int, disable *bool, remark *string) error {
	params := map[string]interface{}{}
	if value != nil {
		params["value"] = *value
	}
	if order != nil {
		params["order"] = *order
	}
	if disable != nil {
		params["disable"] = *disable
	}
	if remark != nil {
		params["remark"] = *remark
	}
	return glb.DB.Model(&domain.Wordcase{}).Where("id = ?", wordcaseID).Updates(params).Error
}

// QueryWordcase implements domain.QueryWordcase
func (u *wordcaseRepo) QueryWordcase(wordcaseID *int64, groupName *string, key *string, limit int, offset int) ([]*domain.Wordcase, error) {
	if wordcaseID != nil {
		var res []*domain.Wordcase
		if err := glb.DB.Where("id = ?", wordcaseID).Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	}
	var total int64
	var res []*domain.Wordcase
	conn := glb.DB.Model(&domain.Wordcase{})

	if groupName != nil {
		conn = conn.Where("group_name = ?", groupName)
	}
	if key != nil {
		conn = conn.Where("key = ?", key)
	}
	if limit != 0 {
		conn.Limit(limit)
	}
	if err := conn.Count(&total).Error; err != nil {
		return nil, err
	}
	conn = conn.Offset(offset)
	if err := conn.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

func (u *wordcaseRepo) MGetWordcases(wordcaseIDs []int64) ([]*domain.Wordcase, error) {
	var res []*domain.Wordcase
	if err := glb.DB.Where("id in (?)", wordcaseIDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
