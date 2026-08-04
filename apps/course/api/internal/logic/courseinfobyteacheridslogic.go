// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseInfoByTeacherIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseInfoByTeacherIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoByTeacherIdsLogic {
	return &CourseInfoByTeacherIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseInfoByTeacherIdsLogic) CourseInfoByTeacherIds(req *types.CourseInfoByTeacherIdsReq) (resp []types.TeacherCourseCountVO, err error) {
	// todo: add your logic here and delete this line

	return
}
