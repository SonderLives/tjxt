// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/apps/media/rpc/internal/config"
	"tjxt/apps/media/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	MediaModel        model.MediaModel
	FileModel         model.FileModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:    c,
		MediaModel: model.NewMediaModel(conn, c.Cache),
		FileModel:  model.NewFileModel(conn, c.Cache),
	}
}