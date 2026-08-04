// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PagePlansLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPagePlansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PagePlansLogic {
	return &PagePlansLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PagePlansLogic) PagePlans(req *types.PageRequest) (resp *types.Result, err error) {
	return NewPageLessonsLogic(l.ctx, l.svcCtx).PageLessons(req)
}
