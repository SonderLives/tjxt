// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	backendlogic "tjxt/service/message/api/backend/logic"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"

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
	r, err := backendlogic.NewListPublicNoticesLogic(l.ctx, l.svcCtx.Backend).ListPublicNotices(&pb.PageReq{
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
