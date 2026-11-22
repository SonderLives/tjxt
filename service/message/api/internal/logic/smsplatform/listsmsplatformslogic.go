// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package smsplatform

import (
	"context"

	backendlogic "tjxt/service/message/api/backend/logic"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSmsPlatformsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSmsPlatformsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSmsPlatformsLogic {
	return &ListSmsPlatformsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSmsPlatformsLogic) ListSmsPlatforms() (resp *types.SmsPlatformListVO, err error) {
	r, err := backendlogic.NewListSmsPlatformsLogic(l.ctx, l.svcCtx.Backend).ListSmsPlatforms(&pb.Empty{})
	if err != nil {
		return nil, err
	}
	list := make([]types.SmsPlatformVO, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.SmsPlatformVO{
			Id:       item.Id,
			Name:     item.Name,
			Code:     item.Code,
			Priority: item.Priority,
			Status:   item.Status,
		})
	}
	return &types.SmsPlatformListVO{List: list}, nil
}
