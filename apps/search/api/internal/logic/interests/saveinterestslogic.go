// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"context"

	"tjxt/apps/search/api/internal/svc"
	"tjxt/apps/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveInterestsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveInterestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveInterestsLogic {
	return &SaveInterestsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveInterestsLogic) SaveInterests(req *types.SaveInterestsReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
