package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLessonGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonGetLogic {
	return &LessonGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 指定课程的学习信息
func (l *LessonGetLogic) LessonGet(in *pb.LessonRequest) (*pb.LearningLessonVO, error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	lesson, err := l.svcCtx.LearningService.GetLesson(l.ctx, userID, in.CourseId)
	if err != nil {
		return nil, err
	}
	return toLessonVO(lesson), nil
}
