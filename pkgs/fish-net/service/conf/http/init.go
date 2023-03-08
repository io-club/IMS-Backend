package http

import (
	"fishnet/domain"
	"fishnet/service/conf/usecase"
)

var _wordcaseUsecase domain.WordcaseUsecase

func init() {
	_wordcaseUsecase = usecase.NewWordcaseUsecase()
}
