// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCataSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCataSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCataSaveLogic {
	return &CourseCataSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseCataSaveLogic) CourseCataSave(req *types.CataSaveReq) (resp *types.NameExistVO, err error) {
	// todo: add your logic here and delete this line

	return
}
