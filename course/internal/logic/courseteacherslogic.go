// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"common/result"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeachersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersLogic {
	return &CourseTeachersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTeachersLogic) CourseTeachers(req *types.CourseTeacherQuery) (resp *result.R, err error) {
	// This handler is for /courses/teachers/:id - need to get ID from path
	// For now, return empty
	return result.OkData([]*types.CourseTeacherVO{}), nil
}
