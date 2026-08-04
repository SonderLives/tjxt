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

// 课程的学习人数
func (l *LessonCountLogic) LessonCount(in *pb.LessonCountRequest) (*pb.LessonCountReply, error) {
	// todo: add your logic here and delete this line

	return &pb.LessonCountReply{}, nil
}
