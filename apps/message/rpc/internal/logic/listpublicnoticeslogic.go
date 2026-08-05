package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPublicNoticesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPublicNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPublicNoticesLogic {
	return &ListPublicNoticesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 公告 分页查询，按 push_time 倒序
func (l *ListPublicNoticesLogic) ListPublicNotices(in *pb.PageReq) (*pb.PublicNoticeListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	total, err := l.svcCtx.PublicNoticeModel.FindCount(l.ctx)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "统计公告失败")
	}
	list, err := l.svcCtx.PublicNoticeModel.FindList(l.ctx, offset, limit)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "分页查询公告失败")
	}
	resp := &pb.PublicNoticeListReply{
		Total: total,
		List:  make([]*pb.PublicNoticeVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, &pb.PublicNoticeVO{
			Id:         item.Id,
			Title:      item.Title,
			Content:    item.Content,
			Type:       int32(item.Type),
			PushTime:   formatTime(item.PushTime),
			ExpireTime: formatTime(item.ExpireTime),
		})
	}
	return resp, nil
}
