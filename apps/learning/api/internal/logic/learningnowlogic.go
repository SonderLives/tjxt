// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LearningNowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLearningNowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LearningNowLogic {
	return &LearningNowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LearningNowLogic) LearningNow() (resp *types.LearningLessonVO, err error) {
	// todo: add your logic here and delete this line

	return
}
