// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseBaseInfoSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseBaseInfoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoSaveLogic {
	return &CourseBaseInfoSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseBaseInfoSaveLogic) CourseBaseInfoSave(req *types.CourseBaseInfoSaveReq) (resp *types.CourseSaveVO, err error) {
	// todo: add your logic here and delete this line

	return
}
