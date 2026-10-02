// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"context"

	"tjxt/apps/media/rpc/internal/config"
	"tjxt/apps/media/rpc/internal/model"
	"tjxt/apps/media/rpc/internal/store"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	MediaModel model.MediaModel
	FileModel  model.FileModel

	// Store 对象存储客户端；nil 表示未配置 MinIO，logic 回退 mock 实现
	Store *store.Store
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	svcCtx := &ServiceContext{
		Config:     c,
		MediaModel: model.NewMediaModel(conn, c.Cache),
		FileModel:  model.NewFileModel(conn, c.Cache),
	}

	if c.MinIO.Endpoint != "" {
		st, err := store.New(context.Background(), store.Config{
			Endpoint:  c.MinIO.Endpoint,
			AccessKey: c.MinIO.AccessKey,
			SecretKey: c.MinIO.SecretKey,
			Bucket:    c.MinIO.Bucket,
			Secure:    c.MinIO.Secure,
		})
		if err != nil {
			// 存储不可用仅告警，服务以 mock 模式继续运行
			logx.Errorf("init object store failed, fallback to mock storage: %v", err)
		} else {
			svcCtx.Store = st
			logx.Infof("object store ready, endpoint=%s bucket=%s", c.MinIO.Endpoint, c.MinIO.Bucket)
		}
	}

	return svcCtx
}
