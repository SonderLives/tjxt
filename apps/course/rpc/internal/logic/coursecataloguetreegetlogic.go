package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueTreeGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueTreeGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueTreeGetLogic {
	return &CourseCatalogueTreeGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCatalogueTreeGet 管理端查询课程目录树（草稿表 course_catalogue_draft）。
// with_practice 为 false 时过滤掉测试（type=3）节点，see（学员视角）暂不影响树结构。
func (l *CourseCatalogueTreeGetLogic) CourseCatalogueTreeGet(in *pb.CourseCatalogueQueryRequest) (*pb.CatalogueTreeList, error) {
	list, err := l.svcCtx.CourseCatalogueDraftModel.ListByCourseId(l.ctx, in.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录草稿失败")
	}
	if !in.WithPractice {
		filtered := make([]*model.CourseCatalogueDraft, 0, len(list))
		for _, c := range list {
			if c.Type == CatalogueTypePractice {
				continue
			}
			filtered = append(filtered, c)
		}
		list = filtered
	}
	return &pb.CatalogueTreeList{Items: buildDraftCataTree(list)}, nil
}

// buildDraftCataTree 由扁平的草稿目录列表构造章-节树，规则同正式表：
// 章（type=1）parent_catalogue_id = 0 作为根，节/测试挂到所属章的 Sections 下。
func buildDraftCataTree(list []*model.CourseCatalogueDraft) []*pb.CourseChapterInfo {
	nodes := make(map[int64]*pb.CourseChapterInfo, len(list))
	ordered := make([]*pb.CourseChapterInfo, 0, len(list))
	for _, c := range list {
		n := &pb.CourseChapterInfo{
			Id:            c.Id,
			Name:          c.Name,
			Index:         int32(c.CIndex),
			Type:          int32(c.Type),
			MediaId:       c.MediaId,
			MediaName:     c.VideoName,
			MediaDuration: c.MediaDuration,
			Trailer:       c.Trailer == 1,
			CanUpdate:     c.CanUpdate == 1,
		}
		nodes[c.Id] = n
		ordered = append(ordered, n)
	}
	roots := make([]*pb.CourseChapterInfo, 0, len(list))
	for i, c := range list {
		n := ordered[i]
		if c.ParentCatalogueId != 0 {
			if parent, ok := nodes[c.ParentCatalogueId]; ok {
				parent.Sections = append(parent.Sections, n)
				continue
			}
		}
		roots = append(roots, n)
	}
	return roots
}
