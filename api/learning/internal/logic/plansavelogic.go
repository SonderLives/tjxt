package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlanSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlanSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanSaveLogic {
	return &PlanSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建/更新学习计划
func (l *PlanSaveLogic) PlanSave(req *types.PlanSaveReq) (*types.OkVO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	if _, err := l.svcCtx.LearningRpc.PlanSave(l.ctx, &pb.PlanSaveRequest{
		CourseId: req.CourseId,
		Freq:     req.Freq,
	}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
