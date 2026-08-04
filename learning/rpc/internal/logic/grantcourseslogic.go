package logic

import (
	"context"

	"learning/rpc/internal/svc"
	"learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantCoursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGrantCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantCoursesLogic {
	return &GrantCoursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GrantCoursesLogic) GrantCourses(in *pb.GrantCoursesRequest) (*pb.EmptyReply, error) {
	if err := l.svcCtx.LessonService.AddUserLessons(l.ctx, in.UserId, in.CourseIds); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
