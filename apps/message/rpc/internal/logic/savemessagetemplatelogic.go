package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/model"
	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMessageTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveMessageTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMessageTemplateLogic {
	return &SaveMessageTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短信模板 新增/更新
func (l *SaveMessageTemplateLogic) SaveMessageTemplate(in *pb.MessageTemplateSaveReq) (*pb.IdReply, error) {
	if in.Name == "" || in.PlatformCode == "" || in.SignName == "" || in.ThirdTemplateCode == "" || in.Content == "" {
		return nil, xerr.BadRequestf("模板名称/平台代号/签名/第三方模板code/内容不能为空")
	}

	// 校验短信平台存在且启用
	platform, err := l.svcCtx.SmsThirdPlatformModel.FindByCode(l.ctx, in.PlatformCode)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.BadRequestf("短信平台不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询短信平台失败")
	}
	if platform.Status != MessageTemplateStatusEnabled {
		return nil, xerr.BadRequestf("短信平台未启用")
	}

	// 校验通知模板存在
	if _, err := l.svcCtx.NoticeTemplateModel.FindOne(l.ctx, in.TemplateId); err != nil {
		if isNotFound(err) {
			return nil, xerr.BadRequestf("通知模板不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询通知模板失败")
	}

	operator := userIdFromCtx(l.ctx)
	if in.Id == 0 {
		data := &model.MessageTemplate{
			Id:                nextID(),
			Name:              in.Name,
			PlatformCode:      in.PlatformCode,
			SignName:          in.SignName,
			ThirdTemplateCode: in.ThirdTemplateCode,
			Content:           in.Content,
			TemplateId:        in.TemplateId,
			Status:            int64(in.Status),
			Creater:           operator,
			Updater:           operator,
		}
		if _, err := l.svcCtx.MessageTemplateModel.Insert(l.ctx, data); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "新增短信模板失败")
		}
		return &pb.IdReply{Id: data.Id}, nil
	}

	old, err := l.svcCtx.MessageTemplateModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("短信模板不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询短信模板失败")
	}
	old.Name = in.Name
	old.PlatformCode = in.PlatformCode
	old.SignName = in.SignName
	old.ThirdTemplateCode = in.ThirdTemplateCode
	old.Content = in.Content
	old.TemplateId = in.TemplateId
	old.Status = int64(in.Status)
	old.Updater = operator
	if err := l.svcCtx.MessageTemplateModel.Update(l.ctx, old); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新短信模板失败")
	}
	return &pb.IdReply{Id: old.Id}, nil
}
