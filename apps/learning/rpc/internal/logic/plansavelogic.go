package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/rpc/internal/svc"
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

// 创建/更新学习计划（每周章节数）
func (l *PlanSaveLogic) PlanSave(in *pb.PlanSaveRequest) (*pb.Empty, error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.LearningService.CreatePlan(l.ctx, userID, in.CourseId, in.Freq); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
