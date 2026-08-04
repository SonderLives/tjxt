// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlanSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlanSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanSaveLogic {
	return &PlanSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlanSaveLogic) PlanSave(req *types.PlanSaveReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
