package logic

import (
	"context"

	"learning/rpc/internal/svc"
	"learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLessonLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLessonLogic {
	return &GetLessonLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLessonLogic) GetLesson(in *pb.LessonRequest) (*pb.LessonReply, error) {
	lesson, err := l.svcCtx.LessonService.GetLesson(l.ctx, in.UserId, in.CourseId)
	if err != nil {
		return nil, err
	}
	return lessonReply(lesson), nil
}
