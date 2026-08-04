// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoursePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoursePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePageLogic {
	return &CoursePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoursePageLogic) CoursePage(req *types.CoursePageReq) (resp *types.CoursePageReply, err error) {
	// todo: add your logic here and delete this line

	return
}
