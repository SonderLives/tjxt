// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"
	messageclient "tjxt/rpc/message/client/message"

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
	r, err := l.svcCtx.MessageRpc.SaveNoticeTask(l.ctx, &messageclient.NoticeTaskSaveReq{
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
