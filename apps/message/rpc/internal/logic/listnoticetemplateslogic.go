package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticeTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNoticeTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeTemplatesLogic {
	return &ListNoticeTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知模板 分页查询，按 update_time 倒序
func (l *ListNoticeTemplatesLogic) ListNoticeTemplates(in *pb.PageReq) (*pb.NoticeTemplateListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	total, err := l.svcCtx.NoticeTemplateModel.FindCount(l.ctx)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "统计通知模板失败")
	}
	list, err := l.svcCtx.NoticeTemplateModel.FindList(l.ctx, offset, limit)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "分页查询通知模板失败")
	}
	resp := &pb.NoticeTemplateListReply{
		Total: total,
		List:  make([]*pb.NoticeTemplateVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, toNoticeTemplateVO(item))
	}
	return resp, nil
}
