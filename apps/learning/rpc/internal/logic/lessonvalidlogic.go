package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
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

// 学员是否已报名（返回 lesson_id；未报名返回 0/NotFound 由调用方判）
func (l *LessonValidLogic) LessonValid(in *pb.LessonRequest) (*pb.LessonValidReply, error) {
	// todo: add your logic here and delete this line

	return &pb.LessonValidReply{}, nil
}
