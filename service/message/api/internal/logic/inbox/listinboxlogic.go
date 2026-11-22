// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inbox

import (
	"context"

	backendlogic "tjxt/service/message/api/backend/logic"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"
	"tjxt/common/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInboxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInboxLogic {
	return &ListInboxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInboxLogic) ListInbox(req *types.InboxListReq) (resp *types.InboxListVO, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	r, err := backendlogic.NewListInboxLogic(l.ctx, l.svcCtx.Backend).ListInbox(&pb.InboxPageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
		UserId:   userId,
		Type:     req.Type,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.UserInboxVO, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.UserInboxVO{
			Id:         item.Id,
			UserId:     item.UserId,
			Type:       item.Type,
			Title:      item.Title,
			Content:    item.Content,
			IsRead:     item.IsRead,
			Publisher:  item.Publisher,
			PushTime:   item.PushTime,
			ExpireTime: item.ExpireTime,
		})
	}
	return &types.InboxListVO{Total: r.Total, List: list}, nil
}
