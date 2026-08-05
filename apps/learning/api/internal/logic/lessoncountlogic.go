package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLessonCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonCountLogic {
	return &LessonCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 课程的学习人数（公开统计，无需登录）
func (l *LessonCountLogic) LessonCount(req *types.LessonRequest) (*types.LessonCountVO, error) {
	reply, err := l.svcCtx.LearningRpc.LessonCount(l.ctx, &pb.LessonCountRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}
	return &types.LessonCountVO{Count: reply.Count}, nil
}
