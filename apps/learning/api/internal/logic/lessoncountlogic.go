// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonCountLogic {
	return &LessonCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonCountLogic) LessonCount(req *types.LessonRequest) (resp *types.LessonCountVO, err error) {
	// todo: add your logic here and delete this line

	return
}
