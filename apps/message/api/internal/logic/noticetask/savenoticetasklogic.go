// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
