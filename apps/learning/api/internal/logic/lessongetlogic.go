package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
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
func (l *LessonGetLogic) LessonGet(req *types.LessonRequest) (*types.LearningLessonVO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	vo, err := l.svcCtx.LearningRpc.LessonGet(l.ctx, &pb.LessonRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}
	enrichLessons(l.ctx, l.svcCtx, []*pb.LearningLessonVO{vo})
	res := toLessonVOTypes(vo)
	return &res, nil
}
