// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"
	mediaclient "tjxt/apps/media/rpc/media"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMediaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaListLogic {
	return &MediaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MediaListLogic) MediaList(req *types.MediaListReq) (resp *types.MediaListVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.MediaRpc.MediaList(l.ctx, &mediaclient.MediaListRequest{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
		Name:     req.Name,
		SortBy:   req.SortBy,
		IsAsc:    req.IsAsc,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.MediaVO, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, types.MediaVO{
			Id:         item.Id,
			Filename:   item.Filename,
			MediaUrl:   item.MediaUrl,
			CoverUrl:   item.CoverUrl,
			Duration:   item.Duration,
			Size:       item.Size,
			Status:     item.Status,
			Creater:    item.Creater,
			CreateTime: item.CreateTime,
			UseTimes:   item.UseTimes,
		})
	}
	return &types.MediaListVO{Total: rpcResp.Total, List: list}, nil
}
