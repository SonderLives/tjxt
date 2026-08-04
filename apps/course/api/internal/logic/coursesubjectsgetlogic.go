// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSubjectsGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSubjectsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsGetLogic {
	return &CourseSubjectsGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseSubjectsGetLogic) CourseSubjectsGet(req *types.IdPathReq) (resp []types.CataSubjectVO, err error) {
	// todo: add your logic here and delete this line

	return
}
