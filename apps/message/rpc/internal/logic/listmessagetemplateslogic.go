package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessageTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessageTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessageTemplatesLogic {
	return &ListMessageTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短信模板 分页查询，按 update_time 倒序
func (l *ListMessageTemplatesLogic) ListMessageTemplates(in *pb.PageReq) (*pb.MessageTemplateListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	total, err := l.svcCtx.MessageTemplateModel.FindCount(l.ctx)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "统计短信模板失败")
	}
	list, err := l.svcCtx.MessageTemplateModel.FindList(l.ctx, offset, limit)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "分页查询短信模板失败")
	}
	resp := &pb.MessageTemplateListReply{
		Total: total,
		List:  make([]*pb.MessageTemplateVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, toMessageTemplateVO(item))
	}
	return resp, nil
}
