package logic

import (
	"context"

	"tjxt/pkg/auth"

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
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.LearningService.CommitRecord(l.ctx, userID, in.LessonId, in.SectionId, in.Moment, in.Duration, in.CommitTime); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
