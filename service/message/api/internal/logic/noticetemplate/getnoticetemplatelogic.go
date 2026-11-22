// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	backendlogic "tjxt/service/message/api/backend/logic"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeTemplateLogic {
	return &GetNoticeTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoticeTemplateLogic) GetNoticeTemplate(req *types.IdPathReq) (resp *types.NoticeTemplateVO, err error) {
	r, err := backendlogic.NewGetNoticeTemplateLogic(l.ctx, l.svcCtx.Backend).GetNoticeTemplate(&pb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.NoticeTemplateVO{
		Id:            r.Id,
		Name:          r.Name,
		Code:          r.Code,
		Type:          r.Type,
		Status:        r.Status,
		Title:         r.Title,
		Content:       r.Content,
		IsSmsTemplate: r.IsSmsTemplate,
		CreateTime:    r.CreateTime,
	}, nil
}
