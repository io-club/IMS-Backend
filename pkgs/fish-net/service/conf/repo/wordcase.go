package repo

import (
	"IMS-Backend/pkgs/fish-net/common/consts"
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/glb"
	"IMS-Backend/pkgs/fish-net/util"
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
	return glb.DB.Where("id = ?", wordcaseId).Delete(&domain.Wordcase{}).Error
}

// UpdateWordcase implements domain.WordcaseRepo
func (*wordcaseRepo) UpdateWordcase(wordcaseID int64, value *string, order *int, disable *bool, remark *string) error {
	params := map[string]interface{}{}
	if value != nil && *value != "" {
		params["value"] = *value
	}
	if order != nil && *order != 0 {
		params["order"] = *order
	}
	if disable != nil {
		params["disable"] = *disable
	}
	if remark != nil && *remark != "" {
		params["remark"] = *remark
	}
	return glb.DB.Model(&domain.Wordcase{}).Where("id = ?", wordcaseID).Updates(params).Error
}

// QueryWordcase implements domain.QueryWordcase
func (u *wordcaseRepo) QueryWordcase(wordcaseID *int64, groupName *string, key *string, limit int, offset int) ([]*domain.Wordcase, error) {
	if wordcaseID != nil && *wordcaseID > 0 {
		glb.LOG.Info(util.SPrettyLog("QueryWordcase", "wordcaseID", wordcaseID))
		var res []*domain.Wordcase
		if err := glb.DB.Where("id = ?", wordcaseID).Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	}
	glb.LOG.Info(util.SPrettyLog("QueryWordcase", "groupName", groupName, "key", key, "limit", limit, "offset", offset))
	var total int64
	var res []*domain.Wordcase
	conn := glb.DB.Model(&domain.Wordcase{})

	if groupName != nil && *groupName != "" {
		conn = conn.Where("group_name = ?", groupName)
	}
	if key != nil && *key != "" {
		conn = conn.Where("key = ?", key)
	}
	if limit == 0 {
		limit = consts.DefaultLimit
	}
	conn = conn.Limit(limit).Offset(offset)
	if err := conn.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := conn.Find(&res).Order("id desc").Error; err != nil {
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
