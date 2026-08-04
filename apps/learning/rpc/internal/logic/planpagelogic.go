package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlanPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPlanPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlanPageLogic {
	return &PlanPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我的学习计划分页
func (l *PlanPageLogic) PlanPage(in *pb.LessonPageRequest) (*pb.PlanPageReply, error) {
	// todo: add your logic here and delete this line

	return &pb.PlanPageReply{}, nil
}
