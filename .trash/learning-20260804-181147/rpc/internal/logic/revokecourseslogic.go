package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb/pb"

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
	if err := l.svcCtx.LessonService.RevokeCourses(l.ctx, in.UserId, in.CourseIds); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
