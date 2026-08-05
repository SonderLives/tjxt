// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
	messageclient "tjxt/apps/message/rpc/message"

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
	r, err := l.svcCtx.MessageRpc.GetNoticeTask(l.ctx, &messageclient.IdReq{Id: req.Id})
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
