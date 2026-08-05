package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogsGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogsGetLogic {
	return &CourseCatalogsGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程目录 (catalogue) =====
// CourseCatalogsGet 学员端课程目录详情：读正式表 course_catalogue 构造章-节树。
func (l *CourseCatalogsGetLogic) CourseCatalogsGet(in *pb.IdRequest) (*pb.CourseAndSectionView, error) {
	course, err := l.svcCtx.CourseModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("课程不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}

	list, err := l.svcCtx.CourseCatalogueModel.ListByCourseId(l.ctx, in.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}

	// 小节数量：目录中 type=2（节）的数量，无目录数据时回落到课程表冗余字段
	var sections int64
	for _, c := range list {
		if c.Type == CatalogueTypeSection {
			sections++
		}
	}
	if sections == 0 {
		sections = formatNullInt64(course.SectionNum)
	}

	// 教师姓名/头像归属 user 服务，course 侧未接线，此处留空；lesson_id 属 learning 服务，填 0。
	return &pb.CourseAndSectionView{
		Id:              course.Id,
		Name:            course.Name,
		CoverUrl:        course.CoverUrl,
		Sections:        sections,
		TeacherIcon:     "",
		TeacherName:     "",
		LessonId:        0,
		LatestSectionId: 0,
		Chapters:        buildCataTree(list),
	}, nil
}

// buildCataTree 由扁平的正式课程目录列表构造章-节树：
// 章（type=1）parent_catalogue_id = 0 作为根，节/测试（type=2/3）挂到所属章的 Sections 下。
// 入参已按 c_index 升序排列，构造过程保持该顺序。
func buildCataTree(list []*model.CourseCatalogue) []*pb.CourseChapterInfo {
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
			CanUpdate:     false,
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
