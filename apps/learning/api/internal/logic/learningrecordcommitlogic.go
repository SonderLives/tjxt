package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
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

// 提交学习记录
func (l *LearningRecordCommitLogic) LearningRecordCommit(req *types.LearningRecordCommitReq) (*types.OkVO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	// section_type 字符串（VIDEO/EXAM）映射为 proto 数字占位（服务层目前忽略该字段）
	var sectionType int32
	switch req.SectionType {
	case "EXAM":
		sectionType = 2
	default:
		sectionType = 1
	}
	if _, err := l.svcCtx.LearningRpc.LearningRecordCommit(l.ctx, &pb.LearningRecordCommitRequest{
		LessonId:    req.LessonId,
		SectionId:   req.SectionId,
		SectionType: sectionType,
		Moment:      req.Moment,
		Duration:    req.Duration,
		CommitTime:  req.CommitTime,
	}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
