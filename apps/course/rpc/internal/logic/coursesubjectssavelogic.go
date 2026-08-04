package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSubjectsSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSubjectsSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsSaveLogic {
	return &CourseSubjectsSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseSubjectsSaveLogic) CourseSubjectsSave(in *pb.CourseSubjectsSaveRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
