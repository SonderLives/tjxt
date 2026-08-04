// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlanPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanPageLogic {
	return &PlanPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlanPageLogic) PlanPage(req *types.PageRequest) (resp *types.PlanPageReply, err error) {
	// todo: add your logic here and delete this line

	return
}
