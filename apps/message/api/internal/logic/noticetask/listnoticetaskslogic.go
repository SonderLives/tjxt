// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
