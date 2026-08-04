// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LearningRecordsByCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLearningRecordsByCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LearningRecordsByCourseLogic {
	return &LearningRecordsByCourseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LearningRecordsByCourseLogic) LearningRecordsByCourse(req *types.LessonRequest) (resp *types.LearningLessonDTO, err error) {
	// todo: add your logic here and delete this line

	return
}
