package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseInfoByTeacherIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseInfoByTeacherIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoByTeacherIdsLogic {
	return &CourseInfoByTeacherIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseInfoByTeacherIdsLogic) CourseInfoByTeacherIds(in *pb.TeacherIdsRequest) (*pb.TeacherCourseCountList, error) {
	// todo: add your logic here and delete this line

	return &pb.TeacherCourseCountList{}, nil
}
