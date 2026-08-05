// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetTodayDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTodayDataLogic {
	return &SetTodayDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetTodayDataLogic) SetTodayData(req *types.TodayDataSetReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
