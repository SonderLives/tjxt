// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/response"
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

func (l *CourseBaseInfoSaveLogic) CourseBaseInfoSave(req *types.CourseBaseInfoSaveDTO) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	vo, err := l.svcCtx.CourseService.SaveBaseInfo(l.ctx, req, userID)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}
