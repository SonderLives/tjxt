package logic

import (
	"context"

	"learning/rpc/internal/svc"
	"learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLessonsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCountLessonsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountLessonsLogic {
	return &CountLessonsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CountLessonsLogic) CountLessons(in *pb.CourseRequest) (*pb.CountReply, error) {
	count, err := l.svcCtx.LessonService.CountLessons(l.ctx, in.CourseId)
	if err != nil {
		return nil, err
	}
	return &pb.CountReply{Count: count}, nil
}
