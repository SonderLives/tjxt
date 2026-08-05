package svc

import (
	"tjxt/apps/remark/rpc/internal/config"
	"tjxt/apps/remark/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config          config.Config
	LikeRecordModel model.LikeRecordModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:          c,
		LikeRecordModel: model.NewLikeRecordModel(conn, c.Cache),
	}
}