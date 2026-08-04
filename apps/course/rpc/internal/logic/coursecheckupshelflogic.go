package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCheckUpShelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCheckUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckUpShelfLogic {
	return &CourseCheckUpShelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCheckUpShelfLogic) CourseCheckUpShelf(in *pb.IdRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
