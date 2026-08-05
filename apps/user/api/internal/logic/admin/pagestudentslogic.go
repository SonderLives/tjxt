// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"context"

	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageStudentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageStudentsLogic {
	return &PageStudentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageStudentsLogic) PageStudents(req *types.PageRequest) (resp *types.UserPageVO, err error) {
	// todo: add your logic here and delete this line

	return
}
