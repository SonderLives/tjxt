package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSearchIndexInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSearchIndexInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSearchIndexInfoListLogic {
	return &CourseSearchIndexInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseSearchIndexInfoList 分页拉取「已上架」课程的索引数据，供 search 服务
// 全量重建 ES 索引。复用 CourseSearchInfoForIndex 的字段映射，仅以 status 过滤
// 上架课程，避免 search 侧逐条回源（N+1）。
func (l *CourseSearchIndexInfoListLogic) CourseSearchIndexInfoList(in *pb.CourseSearchIndexInfoListRequest) (*pb.CourseSearchIndexInfoListReply, error) {
	pageNo := in.PageNo
	if pageNo < 1 {
		pageNo = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 100
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	list, total, err := l.svcCtx.CourseModel.PageQuery(l.ctx, model.CoursePageFilter{
		Status: CourseStatusUpShelf,
	}, (pageNo-1)*pageSize, pageSize)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程索引列表失败")
	}

	items := make([]*pb.CourseSearchIndexInfo, 0, len(list))
	for _, c := range list {
		items = append(items, &pb.CourseSearchIndexInfo{
			Id:            c.Id,
			Name:          c.Name,
			CoverUrl:      c.CoverUrl,
			Price:         c.Price,
			Score:         c.Score,
			Sold:          0, // 销量来自 trade 服务，course 库无对应列
			Sections:      formatNullInt64(c.SectionNum),
			Free:          int32(c.Free),
			CourseType:    int32(c.CourseType),
			Enable:        1, // 已按上架状态过滤，恒为 1
			CategoryIdLv1: c.FirstCateId,
			CategoryIdLv2: c.SecondCateId,
			CategoryIdLv3: c.ThirdCateId,
			CreateTime:    formatTime(c.CreateTime),
			PublishTime:   formatNullTime(c.PublishTime),
			Duration:      formatNullInt64(c.MediaDuration),
		})
	}

	return &pb.CourseSearchIndexInfoListReply{
		Total: total,
		Items: items,
	}, nil
}
