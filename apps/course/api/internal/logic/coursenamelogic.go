// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseNameLogic {
	return &CourseNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseNameLogic) CourseName(req *types.CourseNameReq) (resp []int64, err error) {
	// todo: add your logic here and delete this line

	return
}
