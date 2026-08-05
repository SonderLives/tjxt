// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"context"

	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageTeachersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageTeachersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageTeachersLogic {
	return &PageTeachersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageTeachersLogic) PageTeachers(req *types.PageRequest) (resp *types.UserPageVO, err error) {
	// todo: add your logic here and delete this line

	return
}
