package common

type PageRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
