package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticeTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNoticeTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeTasksLogic {
	return &ListNoticeTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知任务 分页查询，按 update_time 倒序
func (l *ListNoticeTasksLogic) ListNoticeTasks(in *pb.PageReq) (*pb.NoticeTaskListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	total, err := l.svcCtx.NoticeTaskModel.FindCount(l.ctx)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "统计通知任务失败")
	}
	list, err := l.svcCtx.NoticeTaskModel.FindList(l.ctx, offset, limit)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "分页查询通知任务失败")
	}
	resp := &pb.NoticeTaskListReply{
		Total: total,
		List:  make([]*pb.NoticeTaskVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, toNoticeTaskVO(item))
	}
	return resp, nil
}
