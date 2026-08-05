// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPublicNoticesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPublicNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPublicNoticesLogic {
	return &ListPublicNoticesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPublicNoticesLogic) ListPublicNotices(req *types.PageReq) (resp *types.PublicNoticeListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
