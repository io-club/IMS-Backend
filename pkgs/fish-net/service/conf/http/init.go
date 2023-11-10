package http

import (
	"IMS-Backend/pkgs/fish-net/domain"
	"IMS-Backend/pkgs/fish-net/service/conf/usecase"
)

var _wordcaseUsecase domain.WordcaseUsecase

func init() {
	_wordcaseUsecase = usecase.NewWordcaseUsecase()
}
