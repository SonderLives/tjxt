package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &pb.LearningLessonVO{}, nil
}
