// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonValidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonValidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonValidLogic {
	return &LessonValidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonValidLogic) LessonValid(req *types.LessonRequest) (resp *types.LessonValidVO, err error) {
	// todo: add your logic here and delete this line

	return
}
