// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeachersSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersSaveLogic {
	return &CourseTeachersSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTeachersSaveLogic) CourseTeachersSave(req *types.CourseTeacherSaveReq) (resp *types.NameExistVO, err error) {
	// todo: add your logic here and delete this line

	return
}
