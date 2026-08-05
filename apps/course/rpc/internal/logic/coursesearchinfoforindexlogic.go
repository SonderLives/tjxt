package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSearchInfoForIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSearchInfoForIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSearchInfoForIndexLogic {
	return &CourseSearchInfoForIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseSearchInfoForIndex 查询课程用于搜索引擎建索引的字段（正式课程表）。
func (l *CourseSearchInfoForIndexLogic) CourseSearchInfoForIndex(in *pb.IdRequest) (*pb.CourseSearchIndexInfo, error) {
	if in.Id == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	c, err := l.svcCtx.CourseModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("课程不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}

	var enable int32
	if c.Status == CourseStatusUpShelf {
		enable = 1
	}

	return &pb.CourseSearchIndexInfo{
		Id:            c.Id,
		Name:          c.Name,
		CoverUrl:      c.CoverUrl,
		Price:         c.Price,
		Score:         c.Score,
		Sold:          0, // 销量来自 trade 服务，course 库无对应列
		Sections:      formatNullInt64(c.SectionNum),
		Free:          int32(c.Free),
		CourseType:    int32(c.CourseType),
		Enable:        enable,
		CategoryIdLv1: c.FirstCateId,
		CategoryIdLv2: c.SecondCateId,
		CategoryIdLv3: c.ThirdCateId,
		CreateTime:    formatTime(c.CreateTime),
		PublishTime:   formatNullTime(c.PublishTime),
		Duration:      formatNullInt64(c.MediaDuration),
	}, nil
}
