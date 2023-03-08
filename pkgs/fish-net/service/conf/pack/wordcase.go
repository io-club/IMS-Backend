package pack

import (
	"fishnet/domain"
	"fishnet/service/common"
)

type CreateWordcaseRequest struct {
	Group   string `json:"group"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Order   int    `json:"order"`
	Disable bool   `json:"disable"`
	Remark  string `json:"remark"`
}

type UpdateWordcaseRequest struct {
	Value   string `json:"value"`
	Order   int    `json:"order"`
	Disable bool   `json:"disable"`
	Remark  string `json:"remark"`
}

type QueryWordcaseRequest struct {
	Group   string `json:"group"`
	Key     string `json:"value"`
	Order   int    `json:"order"`
	Disable bool   `json:"disable"`
	Remark  string `json:"remark"`
	common.PageRequest
}

type QueryWordcaseResponseEntry map[string]QueryWordcaseResponseValue
type QueryWordcaseResponseGroup map[string]QueryWordcaseResponseEntry

type QueryWordcaseResponseValue struct {
	Value   string `json:"value"`
	Order   int    `json:"order"`
	Disable bool   `json:"disable"`
	Remark  string `json:"remark"`
}

func Keys(wordcases []*domain.Wordcase) []string {
	res := make([]string, len(wordcases))
	for _, wordcase := range wordcases {
		res = append(res, wordcase.Key)
	}
	return res
}

func Value(wordcase *domain.Wordcase) QueryWordcaseResponseValue {
	return QueryWordcaseResponseValue{
		Value:   wordcase.Value,
		Order:   wordcase.Order,
		Disable: wordcase.Disable,
		Remark:  wordcase.Remark,
	}
}

func Values(wordcases []*domain.Wordcase) []QueryWordcaseResponseValue {
	res := make([]QueryWordcaseResponseValue, len(wordcases))
	for _, wordcase := range wordcases {
		res = append(res, Value(wordcase))
	}
	return res
}

func Entries(wordcases []*domain.Wordcase) map[string]QueryWordcaseResponseValue {
	res := make(map[string]QueryWordcaseResponseValue)
	for _, wordcase := range wordcases {
		res[wordcase.Key] = Value(wordcase)
	}
	return res
}

func Groups(wordcases []*domain.Wordcase) QueryWordcaseResponseGroup {
	res := make(QueryWordcaseResponseGroup)
	for _, wordcase := range wordcases {
		res[wordcase.GroupName] = Entries(wordcases)
	}
	return res
}
