package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseUpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseUpLogic {
	return &CourseUpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseUpLogic) CourseUp(in *pb.IdsRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
