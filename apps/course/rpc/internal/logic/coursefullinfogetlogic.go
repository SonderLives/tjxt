package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseFullInfoGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseFullInfoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseFullInfoGetLogic {
	return &CourseFullInfoGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseFullInfoGet 查询正式课程详情，可按需附带目录树与授课老师 id。
func (l *CourseFullInfoGetLogic) CourseFullInfoGet(in *pb.CourseFullInfoGetRequest) (*pb.CourseFullInfo, error) {
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

	info := &pb.CourseFullInfo{
		Id:              c.Id,
		Name:            c.Name,
		CoverUrl:        c.CoverUrl,
		Price:           c.Price,
		FirstCateId:     c.FirstCateId,
		SecondCateId:    c.SecondCateId,
		ThirdCateId:     c.ThirdCateId,
		ValidDuration:   c.ValidDuration,
		PurchaseEndTime: formatTime(c.PurchaseEndTime),
		SectionNum:      formatNullInt64(c.SectionNum),
	}

	if in.WithCatalogue {
		catalogues, cerr := l.svcCtx.CourseCatalogueModel.ListByCourseId(l.ctx, c.Id)
		if cerr != nil {
			return nil, xerr.Wrap(cerr, xerr.CodeInternal, "查询课程目录失败")
		}
		info.Chapters = l.buildChapterTree(catalogues)
	}

	if in.WithTeachers {
		teachers, terr := l.svcCtx.CourseTeacherModel.ListByCourseId(l.ctx, c.Id)
		if terr != nil {
			return nil, xerr.Wrap(terr, xerr.CodeInternal, "查询课程老师失败")
		}
		ids := make([]int64, 0, len(teachers))
		for _, t := range teachers {
			ids = append(ids, t.TeacherId)
		}
		info.TeacherIds = ids
	}
	return info, nil
}

// buildChapterTree 按 parent_catalogue_id 组织目录树：章 parent=0，节/测试挂到所属章下。
func (l *CourseFullInfoGetLogic) buildChapterTree(list []*model.CourseCatalogue) []*pb.CourseChapterInfo {
	chapters := make([]*pb.CourseChapterInfo, 0, len(list))
	chapterMap := make(map[int64]*pb.CourseChapterInfo, len(list))
	for _, c := range list {
		if c.ParentCatalogueId != 0 {
			continue
		}
		node := l.toChapterNode(c)
		chapterMap[c.Id] = node
		chapters = append(chapters, node)
	}
	for _, c := range list {
		if c.ParentCatalogueId == 0 {
			continue
		}
		node := l.toChapterNode(c)
		if parent, ok := chapterMap[c.ParentCatalogueId]; ok {
			parent.Sections = append(parent.Sections, node)
			continue
		}
		// 找不到所属章时挂到顶层，避免数据丢失
		chapters = append(chapters, node)
	}
	return chapters
}

// toChapterNode 课程目录实体转 pb 节点。
func (l *CourseFullInfoGetLogic) toChapterNode(c *model.CourseCatalogue) *pb.CourseChapterInfo {
	return &pb.CourseChapterInfo{
		Id:            c.Id,
		Name:          c.Name,
		Index:         int32(c.CIndex),
		Type:          int32(c.Type),
		MediaId:       c.MediaId,
		MediaName:     c.VideoName,
		MediaDuration: c.MediaDuration,
		Trailer:       c.Trailer == 1,
	}
}
