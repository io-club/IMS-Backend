package usecase

import (
	"errors"
	"fishnet/domain"
	"fishnet/service/conf/repo"
)

var _wordcaseUsecase domain.WordcaseUsecase

type wordcaseUsecase struct {
	r domain.WordcaseRepo
}

func NewWordcaseUsecase() domain.WordcaseUsecase {
	if _wordcaseUsecase == nil {
		_wordcaseUsecase = &wordcaseUsecase{
			r: repo.NewWordcaseRepo(),
		}
	}
	return _wordcaseUsecase
}

// CreateWordcase implements domain.WordcaseUsecase
func (u *wordcaseUsecase) CreateWordcase(wordcase []*domain.Wordcase) error {
	return u.r.CreateWordcase(wordcase)
}

// DeleteWordcase implements domain.WordcaseUsecase
func (u *wordcaseUsecase) DeleteWordcase(wordcaseID int64) error {
	return u.r.DeleteWordcase(wordcaseID)
}

// MGetWordcases implements domain.WordcaseUsecase
func (u *wordcaseUsecase) MGetWordcases(wordcaseIDs []int64) ([]*domain.Wordcase, error) {
	return u.r.MGetWordcases(wordcaseIDs)
}

// QueryWordcase implements domain.WordcaseUsecase
func (u *wordcaseUsecase) QueryWordcase(wordcaseID *int64, groupName, key *string, limit, offset int) ([]*domain.Wordcase, error) {
	return u.r.QueryWordcase(wordcaseID, groupName, nil, limit, offset)
}

// UpdateWordcase implements domain.WordcaseUsecase
func (u *wordcaseUsecase) UpdateWordcase(wordcaseID int64, value *string, order *int, disable *bool, remark *string) error {
	wordcases, err := u.r.QueryWordcase(&wordcaseID, nil, nil, 1, 0)
	if err != nil {
		return err
	}
	if len(wordcases) == 0 {
		return errors.New("wordcase not found")
		// return domain.ErrWordcaseNotFound
	}
	return u.r.UpdateWordcase(wordcaseID, value, order, disable, remark)
}
