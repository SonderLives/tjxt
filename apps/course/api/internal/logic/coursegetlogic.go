// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGetLogic {
	return &CourseGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseGetLogic) CourseGet(req *types.CourseFullInfoQueryReq) (resp *types.CourseFullInfoDTO, err error) {
	// todo: add your logic here and delete this line

	return
}
