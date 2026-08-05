package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonValidLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLessonValidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonValidLogic {
	return &LessonValidLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 学员是否已报名（返回 lesson_id；未报名时服务层返回 NotFound）
func (l *LessonValidLogic) LessonValid(in *pb.LessonRequest) (*pb.LessonValidReply, error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	lessonID, err := l.svcCtx.LearningService.ValidateLesson(l.ctx, userID, in.CourseId)
	if err != nil {
		return nil, err
	}
	return &pb.LessonValidReply{LessonId: lessonID}, nil
}
