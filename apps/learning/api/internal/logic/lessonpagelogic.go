// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LessonPageLogic {
	return &LessonPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonPageLogic) LessonPage(req *types.PageRequest) (resp *types.LessonPageReply, err error) {
	// todo: add your logic here and delete this line

	return
}
