// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCheckNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCheckNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckNameLogic {
	return &CourseCheckNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseCheckNameLogic) CourseCheckName(req *types.CourseCheckNameReq) (resp *types.NameExistVO, err error) {
	// todo: add your logic here and delete this line

	return
}
