// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"context"

	"tjxt/apps/search/api/internal/svc"
	"tjxt/apps/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInterestsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInterestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInterestsLogic {
	return &GetInterestsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInterestsLogic) GetInterests() (resp *types.InterestsVO, err error) {
	// todo: add your logic here and delete this line

	return
}
