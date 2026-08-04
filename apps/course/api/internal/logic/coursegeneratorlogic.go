// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseGeneratorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseGeneratorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGeneratorLogic {
	return &CourseGeneratorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseGeneratorLogic) CourseGenerator() (resp *types.CourseCataIdVO, err error) {
	// todo: add your logic here and delete this line

	return
}
