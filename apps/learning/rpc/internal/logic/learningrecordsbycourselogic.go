package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LearningRecordsByCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLearningRecordsByCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LearningRecordsByCourseLogic {
	return &LearningRecordsByCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询某课程的学习记录
func (l *LearningRecordsByCourseLogic) LearningRecordsByCourse(in *pb.LessonRequest) (*pb.LearningRecordsReply, error) {
	// todo: add your logic here and delete this line

	return &pb.LearningRecordsReply{}, nil
}
