// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogsGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCatalogsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogsGetLogic {
	return &CourseCatalogsGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseCatalogsGet 查询课程目录（学员端播放页：课程信息 + 章节树）。
func (l *CourseCatalogsGetLogic) CourseCatalogsGet(req *types.IdPathReq) (resp *types.CourseAndSectionVO, err error) {
	view, gerr := l.svcCtx.CourseRpc.CourseCatalogsGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程目录失败")
	}
	resp = &types.CourseAndSectionVO{
		Id:              view.Id,
		Name:            view.Name,
		CoverUrl:        view.CoverUrl,
		Sections:        view.Sections,
		TeacherIcon:     view.TeacherIcon,
		TeacherName:     view.TeacherName,
		LessonId:        view.LessonId,
		LatestSectionId: view.LatestSectionId,
		Chapters:        make([]types.ChapterVO, 0, len(view.Chapters)),
	}
	for _, c := range view.Chapters {
		resp.Chapters = append(resp.Chapters, catalogsGetToChapterVO(c))
	}
	return resp, nil
}

// catalogsGetToChapterVO pb.CourseChapterInfo -> API ChapterVO（递归映射子节）。
func catalogsGetToChapterVO(c *pb.CourseChapterInfo) types.ChapterVO {
	if c == nil {
		return types.ChapterVO{}
	}
	vo := types.ChapterVO{
		Id:            c.Id,
		Name:          c.Name,
		Index:         int64(c.Index),
		Type:          int64(c.Type),
		MediaId:       c.MediaId,
		MediaName:     c.MediaName,
		MediaDuration: c.MediaDuration,
		Trailer:       c.Trailer,
	}
	for _, s := range c.Sections {
		vo.Sections = append(vo.Sections, catalogsGetToChapterVO(s))
	}
	return vo
}
