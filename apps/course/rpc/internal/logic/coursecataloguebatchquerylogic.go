package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueBatchQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueBatchQueryLogic {
	return &CourseCatalogueBatchQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCatalogueBatchQuery 按目录 id 批量查询目录简要信息（正式表）。
func (l *CourseCatalogueBatchQueryLogic) CourseCatalogueBatchQuery(in *pb.IdsRequest) (*pb.CataSimpleList, error) {
	if len(in.Ids) == 0 {
		return &pb.CataSimpleList{Items: []*pb.CataSimple{}}, nil
	}
	list, err := l.svcCtx.CourseCatalogueModel.ListByIdIn(l.ctx, in.Ids)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "批量查询课程目录失败")
	}
	items := make([]*pb.CataSimple, 0, len(list))
	for _, c := range list {
		items = append(items, &pb.CataSimple{
			Id:     c.Id,
			Name:   c.Name,
			Index:  int32(c.CIndex),
			CIndex: int32(c.CIndex),
		})
	}
	return &pb.CataSimpleList{Items: items}, nil
}
