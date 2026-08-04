// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseDownShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDownShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDownShelfLogic {
	return &CourseDownShelfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseDownShelfLogic) CourseDownShelf(req *types.IdPathReq) (resp *types.NameExistVO, err error) {
	// todo: add your logic here and delete this line

	return
}
