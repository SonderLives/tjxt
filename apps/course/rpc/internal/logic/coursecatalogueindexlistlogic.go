package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueIndexListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueIndexListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueIndexListLogic {
	return &CourseCatalogueIndexListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCatalogueIndexList 查询课程目录草稿的序号列表：
// index 为节点在同级中的序号，chapter_index 为所属章的序号（章自身即其自身序号）。
func (l *CourseCatalogueIndexListLogic) CourseCatalogueIndexList(in *pb.IdRequest) (*pb.CataSimpleList, error) {
	list, err := l.svcCtx.CourseCatalogueDraftModel.ListByCourseId(l.ctx, in.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录草稿失败")
	}

	// 先按 c_index 顺序给章编号
	chapterNo := make(map[int64]int32, len(list))
	var no int32
	for _, c := range list {
		if c.Type == CatalogueTypeChapter {
			no++
			chapterNo[c.Id] = no
		}
	}

	// 同级序号：按 parent 分组累加
	levelNo := make(map[int64]int32, len(list))
	items := make([]*pb.CataSimple, 0, len(list))
	for _, c := range list {
		levelNo[c.ParentCatalogueId]++
		chapterIndex := chapterNo[c.Id]
		if c.Type != CatalogueTypeChapter {
			chapterIndex = chapterNo[c.ParentCatalogueId]
		}
		items = append(items, &pb.CataSimple{
			Id:           c.Id,
			Name:         c.Name,
			Index:        levelNo[c.ParentCatalogueId],
			ChapterIndex: chapterIndex,
			CIndex:       int32(c.CIndex),
		})
	}
	return &pb.CataSimpleList{Items: items}, nil
}
