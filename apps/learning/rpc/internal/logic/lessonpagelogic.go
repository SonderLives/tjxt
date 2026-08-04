package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLessonPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonPageLogic {
	return &LessonPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 我的课表分页
func (l *LessonPageLogic) LessonPage(in *pb.LessonPageRequest) (*pb.LessonPageReply, error) {
	// todo: add your logic here and delete this line

	return &pb.LessonPageReply{}, nil
}
