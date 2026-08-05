// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
