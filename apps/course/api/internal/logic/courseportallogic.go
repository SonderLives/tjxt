// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoursePortalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoursePortalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePortalLogic {
	return &CoursePortalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoursePortalLogic) CoursePortal(req *types.CoursePortalReq) (resp *types.CoursePageReply, err error) {
	// todo: add your logic here and delete this line

	return
}
