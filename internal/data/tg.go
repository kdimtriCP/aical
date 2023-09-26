package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kdimtricp/aical/internal/biz"
)

type TGRepo struct {
	data *Data
	log  *log.Helper
}

func NewTGRepo(data *Data, logger log.Logger) biz.TGRepo {
	return &TGRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
