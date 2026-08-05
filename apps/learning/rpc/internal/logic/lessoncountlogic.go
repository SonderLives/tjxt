package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
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
func (l *LessonCountLogic) LessonCount(in *pb.LessonCountRequest) (*pb.LessonCountReply, error) {
	count, err := l.svcCtx.LearningService.CountLessons(l.ctx, in.CourseId)
	if err != nil {
		return nil, err
	}
	return &pb.LessonCountReply{Count: count}, nil
}
