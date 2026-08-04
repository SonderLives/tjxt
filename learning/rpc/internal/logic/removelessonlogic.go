package logic

import (
	"context"

	"learning/rpc/internal/svc"
	"learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveLessonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveLessonLogic {
	return &RemoveLessonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveLessonLogic) RemoveLesson(in *pb.LessonRequest) (*pb.EmptyReply, error) {
	if err := l.svcCtx.LessonService.DeleteCourseFromLesson(l.ctx, in.UserId, in.CourseId); err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
