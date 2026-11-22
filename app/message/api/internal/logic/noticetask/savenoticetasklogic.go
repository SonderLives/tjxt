// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	backendlogic "tjxt/app/message/api/backend/logic"
	"tjxt/app/message/api/backend/pb"
	"tjxt/app/message/api/internal/svc"
	"tjxt/app/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoticeTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoticeTaskLogic {
	return &SaveNoticeTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveNoticeTaskLogic) SaveNoticeTask(req *types.NoticeTaskSaveReq) (resp *types.IdVO, err error) {
	r, err := backendlogic.NewSaveNoticeTaskLogic(l.ctx, l.svcCtx.Backend).SaveNoticeTask(&pb.NoticeTaskSaveReq{
		Id:         req.Id,
		TemplateId: req.TemplateId,
		Name:       req.Name,
		Partial:    req.Partial,
		PushTime:   req.PushTime,
		Interval:   req.Interval,
		ExpireTime: req.ExpireTime,
		MaxTimes:   req.MaxTimes,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: r.Id}, nil
}
