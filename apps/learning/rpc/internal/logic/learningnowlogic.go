package logic

import (
	"context"

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

// 当前正在学习的课程
func (l *LearningNowLogic) LearningNow(in *pb.Empty) (*pb.LearningLessonVO, error) {
	// todo: add your logic here and delete this line

	return &pb.LearningLessonVO{}, nil
}
