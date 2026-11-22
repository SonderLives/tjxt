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

type ListNoticeTasksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNoticeTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeTasksLogic {
	return &ListNoticeTasksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNoticeTasksLogic) ListNoticeTasks(req *types.PageReq) (resp *types.NoticeTaskListVO, err error) {
	r, err := backendlogic.NewListNoticeTasksLogic(l.ctx, l.svcCtx.Backend).ListNoticeTasks(&pb.PageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.NoticeTaskVO, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.NoticeTaskVO{
			Id:         item.Id,
			TemplateId: item.TemplateId,
			Name:       item.Name,
			Partial:    item.Partial,
			PushTime:   item.PushTime,
			Interval:   item.Interval,
			ExpireTime: item.ExpireTime,
			MaxTimes:   item.MaxTimes,
			Finished:   item.Finished,
			CreateTime: item.CreateTime,
		})
	}
	return &types.NoticeTaskListVO{Total: r.Total, List: list}, nil
}
