// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	backendlogic "tjxt/api/message/backend/logic"
	"tjxt/api/message/backend/pb"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeTaskLogic {
	return &GetNoticeTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoticeTaskLogic) GetNoticeTask(req *types.IdPathReq) (resp *types.NoticeTaskVO, err error) {
	r, err := backendlogic.NewGetNoticeTaskLogic(l.ctx, l.svcCtx.Backend).GetNoticeTask(&pb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.NoticeTaskVO{
		Id:         r.Id,
		TemplateId: r.TemplateId,
		Name:       r.Name,
		Partial:    r.Partial,
		PushTime:   r.PushTime,
		Interval:   r.Interval,
		ExpireTime: r.ExpireTime,
		MaxTimes:   r.MaxTimes,
		Finished:   r.Finished,
		CreateTime: r.CreateTime,
	}, nil
}
