// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"context"

	backendlogic "tjxt/service/search/api/backend/logic"
	backendpb "tjxt/service/search/api/backend/pb"
	"tjxt/service/search/api/internal/svc"
	"tjxt/service/search/api/internal/types"
	"tjxt/common/auth"

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

	vo, err := backendlogic.NewGetInterestsLogic(l.ctx, l.svcCtx.Backend).GetInterests(&backendpb.IdReq{Id: userId})
	if err != nil {
		return nil, err
	}
	return &types.InterestsVO{
		Id:         vo.Id,
		Interests:  vo.Interests,
		CreateTime: vo.CreateTime,
		UpdateTime: vo.UpdateTime,
	}, nil
}
