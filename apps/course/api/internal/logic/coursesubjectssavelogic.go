// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSubjectsSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSubjectsSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsSaveLogic {
	return &CourseSubjectsSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseSubjectsSaveLogic) CourseSubjectsSave(req *types.CourseSubjectsSaveReq) (resp *types.NameExistVO, err error) {
	// todo: add your logic here and delete this line

	return
}
