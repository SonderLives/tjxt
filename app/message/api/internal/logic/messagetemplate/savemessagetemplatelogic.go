// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"context"

	backendlogic "tjxt/app/message/api/backend/logic"
	"tjxt/app/message/api/backend/pb"
	"tjxt/app/message/api/internal/svc"
	"tjxt/app/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMessageTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveMessageTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMessageTemplateLogic {
	return &SaveMessageTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveMessageTemplateLogic) SaveMessageTemplate(req *types.MessageTemplateSaveReq) (resp *types.IdVO, err error) {
	r, err := backendlogic.NewSaveMessageTemplateLogic(l.ctx, l.svcCtx.Backend).SaveMessageTemplate(&pb.MessageTemplateSaveReq{
		Id:                req.Id,
		Name:              req.Name,
		PlatformCode:      req.PlatformCode,
		SignName:          req.SignName,
		ThirdTemplateCode: req.ThirdTemplateCode,
		Content:           req.Content,
		TemplateId:        req.TemplateId,
		Status:            req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: r.Id}, nil
}
