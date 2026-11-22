package svc

import (
	"tjxt/apps/message/rpc/internal/config"
	"tjxt/apps/message/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config                config.Config
	MessageTemplateModel  model.MessageTemplateModel
	NoticeTaskModel       model.NoticeTaskModel
	NoticeTemplateModel   model.NoticeTemplateModel
	PublicNoticeModel     model.PublicNoticeModel
	SmsThirdPlatformModel model.SmsThirdPlatformModel
	UserInboxModel        model.UserInboxModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:                c,
		MessageTemplateModel:  model.NewMessageTemplateModel(conn, c.Cache),
		NoticeTaskModel:       model.NewNoticeTaskModel(conn, c.Cache),
		NoticeTemplateModel:   model.NewNoticeTemplateModel(conn, c.Cache),
		PublicNoticeModel:     model.NewPublicNoticeModel(conn, c.Cache),
		SmsThirdPlatformModel: model.NewSmsThirdPlatformModel(conn, c.Cache),
		UserInboxModel:        model.NewUserInboxModel(conn, c.Cache),
	}
}
