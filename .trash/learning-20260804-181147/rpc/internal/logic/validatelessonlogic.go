package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateLessonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateLessonLogic {
	return &ValidateLessonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateLessonLogic) ValidateLesson(in *pb.LessonRequest) (*pb.ValidateReply, error) {
	lesson, err := l.svcCtx.LessonService.GetLesson(l.ctx, in.UserId, in.CourseId)
	if err != nil || lesson.Status == 3 {
		return &pb.ValidateReply{}, nil
	}
	return &pb.ValidateReply{LessonId: lesson.Id}, nil
}
