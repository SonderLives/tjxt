// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeachersGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersGetLogic {
	return &CourseTeachersGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTeachersGetLogic) CourseTeachersGet(req *types.CourseCataQueryReq) (resp []types.TeacherCourseInfoVO, err error) {
	// todo: add your logic here and delete this line

	return
}
