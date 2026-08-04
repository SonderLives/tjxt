// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonGetLogic {
	return &LessonGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonGetLogic) LessonGet(req *types.LessonRequest) (resp *types.LearningLessonVO, err error) {
	// todo: add your logic here and delete this line

	return
}
