package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LearningNowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLearningNowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LearningNowLogic {
	return &LearningNowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 当前正在学习的课程（最近学习过的一条）
func (l *LearningNowLogic) LearningNow(in *pb.Empty) (*pb.LearningLessonVO, error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	lesson, err := l.svcCtx.LearningService.CurrentLesson(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	return toLessonVO(lesson), nil
}
