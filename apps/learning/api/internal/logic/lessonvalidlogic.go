package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
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

// 学员是否已报名（返回 lessonId；未报名时 RPC 返回 NotFound）
func (l *LessonValidLogic) LessonValid(req *types.LessonRequest) (*types.LessonValidVO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	reply, err := l.svcCtx.LearningRpc.LessonValid(l.ctx, &pb.LessonRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}
	return &types.LessonValidVO{LessonId: reply.LessonId}, nil
}
