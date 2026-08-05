package logic

import (
	"context"

	"tjxt/pkg/auth"

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

// 查询某课程的学习记录。
// 注：学习记录不再单独建表，仅以 learning_lesson 的最新进度表示，
// 故历史 records 列表为空，只返回当前最新进度。
func (l *LearningRecordsByCourseLogic) LearningRecordsByCourse(in *pb.LessonRequest) (*pb.LearningRecordsReply, error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	lesson, err := l.svcCtx.LearningService.GetLesson(l.ctx, userID, in.CourseId)
	if err != nil {
		return nil, err
	}
	return &pb.LearningRecordsReply{
		Id:              lesson.Id,
		LatestSectionId: nullInt64(lesson.LatestSectionId),
		Records:         []*pb.LearningRecordDTO{},
	}, nil
}
