package logic

import (
	"context"

	"learning/rpc/internal/svc"
	"learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeCoursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeCoursesLogic {
	return &RevokeCoursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RevokeCoursesLogic) RevokeCourses(in *pb.GrantCoursesRequest) (*pb.EmptyReply, error) {
	for _, courseID := range in.CourseIds {
		if err := l.svcCtx.LessonService.DeleteCourseFromLesson(l.ctx, in.UserId, courseID); err != nil {
			return nil, err
		}
	}
	return &pb.EmptyReply{}, nil
}
