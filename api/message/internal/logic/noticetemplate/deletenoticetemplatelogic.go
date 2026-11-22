// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	backendlogic "tjxt/api/message/backend/logic"
	"tjxt/api/message/backend/pb"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoticeTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoticeTemplateLogic {
	return &DeleteNoticeTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNoticeTemplateLogic) DeleteNoticeTemplate(req *types.IdPathReq) (resp *types.OkVO, err error) {
	if _, err := backendlogic.NewDeleteNoticeTemplateLogic(l.ctx, l.svcCtx.Backend).DeleteNoticeTemplate(&pb.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
