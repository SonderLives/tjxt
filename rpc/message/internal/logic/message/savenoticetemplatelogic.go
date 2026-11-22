package messagelogic

import (
	"context"
	"database/sql"

	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"
	"tjxt/rpc/message/internal/model"
	"tjxt/rpc/message/internal/svc"
	"tjxt/rpc/message/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoticeTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoticeTemplateLogic {
	return &SaveNoticeTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知模板 新增/更新
func (l *SaveNoticeTemplateLogic) SaveNoticeTemplate(in *pb.NoticeTemplateSaveReq) (*pb.IdReply, error) {
	if in.Name == "" || in.Code == "" || in.Content == "" {
		return nil, errorx.BadRequestf("模板名称/代号/内容不能为空")
	}
	operator, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	title := sql.NullString{String: in.Title, Valid: in.Title != ""}
	isSmsTemplate := byte(0)
	if in.IsSmsTemplate {
		isSmsTemplate = 1
	}

	if in.Id == 0 {
		data := &model.NoticeTemplate{
			Id:            nextID(),
			Name:          in.Name,
			Code:          in.Code,
			Type:          int64(in.Type),
			Status:        NoticeTemplateStatusDraft,
			Title:         title,
			Content:       in.Content,
			IsSmsTemplate: isSmsTemplate,
			Creater:       operator,
			Updater:       operator,
		}
		if _, err := l.svcCtx.NoticeTemplateModel.Insert(l.ctx, data); err != nil {
			return nil, errorx.Wrapf(err, errorx.CodeInternal, "新增通知模板失败")
		}
		return &pb.IdReply{Id: data.Id}, nil
	}

	old, err := l.svcCtx.NoticeTemplateModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("通知模板不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询通知模板失败")
	}
	old.Name = in.Name
	old.Code = in.Code
	old.Type = int64(in.Type)
	old.Title = title
	old.Content = in.Content
	old.IsSmsTemplate = isSmsTemplate
	old.Updater = operator
	if err := l.svcCtx.NoticeTemplateModel.Update(l.ctx, old); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新通知模板失败")
	}
	return &pb.IdReply{Id: old.Id}, nil
}
