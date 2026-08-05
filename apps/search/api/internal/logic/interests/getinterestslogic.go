// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"context"

	"tjxt/apps/search/api/internal/svc"
	"tjxt/apps/search/api/internal/types"
	searchclient "tjxt/apps/search/rpc/client/search"
	"tjxt/pkg/auth"

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
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.SearchRpc.GetInterests(l.ctx, &searchclient.IdReq{Id: userId})
	if err != nil {
		return nil, err
	}
	return &types.InterestsVO{
		Id:         rpcResp.Id,
		Interests:  rpcResp.Interests,
		CreateTime: rpcResp.CreateTime,
		UpdateTime: rpcResp.UpdateTime,
	}, nil
}
