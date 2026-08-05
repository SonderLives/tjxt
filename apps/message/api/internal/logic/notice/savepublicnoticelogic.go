// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePublicNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePublicNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePublicNoticeLogic {
	return &SavePublicNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePublicNoticeLogic) SavePublicNotice(req *types.PublicNoticeSaveReq) (resp *types.IdVO, err error) {
	// todo: add your logic here and delete this line

	return
}
