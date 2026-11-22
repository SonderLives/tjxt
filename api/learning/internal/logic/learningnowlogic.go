package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
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
func (l *LearningNowLogic) LearningNow() (*types.LearningLessonVO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	vo, err := l.svcCtx.LearningRpc.LearningNow(l.ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	enrichLessons(l.ctx, l.svcCtx, []*pb.LearningLessonVO{vo})
	res := toLessonVOTypes(vo)
	return &res, nil
}
