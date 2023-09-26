package service

import (
	"github.com/go-kratos/kratos/v2/log"
)

type TGService struct {
	log *log.Helper
}

func NewTGService(logger log.Logger) *TGService {

	return &TGService{
		log: log.NewHelper(log.With(logger, "module", "service/tg")),
	}
}
