package medialogic

import (
	"context"

	"tjxt/rpc/media/internal/svc"
	"tjxt/rpc/media/pb"
	"tjxt/pkg/page"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaListLogic {
	return &MediaListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MediaListLogic) MediaList(in *pb.MediaListRequest) (*pb.MediaListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))
	list, err := l.svcCtx.MediaModel.FindPage(l.ctx, in.Name, in.SortBy, in.IsAsc, offset, limit)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "分页查询媒资失败")
	}
	total, err := l.svcCtx.MediaModel.Count(l.ctx, in.Name)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "统计媒资数量失败")
	}
	resp := &pb.MediaListReply{
		Total: total,
		List:  make([]*pb.MediaVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, toMediaVO(item))
	}
	return resp, nil
}
