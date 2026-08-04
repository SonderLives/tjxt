// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"learning/internal/svc"
	"learning/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LearningLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLearningLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LearningLogic {
	return &LearningLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LearningLogic) Learning() (*types.Result, error) { return success(nil), nil }
