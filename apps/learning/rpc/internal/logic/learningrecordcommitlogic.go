package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LearningRecordCommitLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLearningRecordCommitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LearningRecordCommitLogic {
	return &LearningRecordCommitLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 提交学习记录（更新 lesson 的 latest_section/learn_time/learned_sections）
func (l *LearningRecordCommitLogic) LearningRecordCommit(in *pb.LearningRecordCommitRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
