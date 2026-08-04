package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseMediaSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaSaveLogic {
	return &CourseMediaSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程媒体 =====
func (l *CourseMediaSaveLogic) CourseMediaSave(in *pb.CourseMediaSaveRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
