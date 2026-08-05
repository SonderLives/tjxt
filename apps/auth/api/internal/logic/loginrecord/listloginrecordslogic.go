// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package loginrecord

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoginRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginRecordsLogic {
	return &ListLoginRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoginRecordsLogic) ListLoginRecords(req *types.LoginRecordListReq) (resp *types.LoginRecordListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
