package logic

import (
	"context"

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

// 创建/更新学习计划
func (l *PlanSaveLogic) PlanSave(in *pb.PlanSaveRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
