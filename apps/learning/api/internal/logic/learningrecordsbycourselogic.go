package logic

import (
	"context"

	"tjxt/pkg/auth"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
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

// 查询某课程的学习记录（含最新进度与历史列表）
func (l *LearningRecordsByCourseLogic) LearningRecordsByCourse(req *types.LessonRequest) (*types.LearningLessonDTO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	reply, err := l.svcCtx.LearningRpc.LearningRecordsByCourse(l.ctx, &pb.LessonRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}
	records := make([]types.LearningRecordDTO, 0, len(reply.Records))
	for _, r := range reply.Records {
		records = append(records, types.LearningRecordDTO{
			SectionId: r.SectionId,
			Moment:    r.Moment,
			Duration:  r.Duration,
			Finished:  r.Finished,
		})
	}
	return &types.LearningLessonDTO{
		Id:              reply.Id,
		LatestSectionId: reply.LatestSectionId,
		Records:         records,
	}, nil
}
