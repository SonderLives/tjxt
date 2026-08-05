// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
	messageclient "tjxt/apps/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPublicNoticesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPublicNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPublicNoticesLogic {
	return &ListPublicNoticesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPublicNoticesLogic) ListPublicNotices(req *types.PageReq) (resp *types.PublicNoticeListVO, err error) {
	r, err := l.svcCtx.MessageRpc.ListPublicNotices(l.ctx, &messageclient.PageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.PublicNoticeVO, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.PublicNoticeVO{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			Type:       item.Type,
			PushTime:   item.PushTime,
			ExpireTime: item.ExpireTime,
		})
	}
	return &types.PublicNoticeListVO{Total: r.Total, List: list}, nil
}
