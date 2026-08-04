// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSimpleInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSimpleInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSimpleInfoLogic {
	return &CourseSimpleInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseSimpleInfoLogic) CourseSimpleInfo(req *types.CourseSimpleInfoQueryReq) (resp []types.CourseSimpleInfoDTO, err error) {
	// todo: add your logic here and delete this line

	return
}
