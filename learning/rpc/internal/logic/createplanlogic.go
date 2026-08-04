package logic

import (
	"context"

	"learning/rpc/internal/svc"
	"learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePlanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePlanLogic {
	return &CreatePlanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePlanLogic) CreatePlan(in *pb.CreatePlanRequest) (*pb.EmptyReply, error) {
	if err := l.svcCtx.LessonService.CreatePlan(l.ctx, in.UserId, in.CourseId, int64(in.Freq)); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
