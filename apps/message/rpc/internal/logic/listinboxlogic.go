package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInboxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInboxLogic {
	return &ListInboxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 站内信 分页查询（仅未过期，按 push_time 倒序）
func (l *ListInboxLogic) ListInbox(in *pb.InboxPageReq) (*pb.InboxListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	total, err := l.svcCtx.UserInboxModel.FindCount(l.ctx, in.UserId, int64(in.Type))
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "统计站内信失败")
	}
	list, err := l.svcCtx.UserInboxModel.FindList(l.ctx, in.UserId, int64(in.Type), offset, limit)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "分页查询站内信失败")
	}
	resp := &pb.InboxListReply{
		Total: total,
		List:  make([]*pb.UserInboxVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, toUserInboxVO(item))
	}
	return resp, nil
}
