package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDeleteLogic {
	return &CourseDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseDelete 删除课程：级联清理草稿表与正式表的目录、老师、题目绑定与内容。
func (l *CourseDeleteLogic) CourseDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if in.Id == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	id := in.Id

	// ===== 草稿表 =====
	if err := l.svcCtx.CourseCataSubjectDraftModel.DeleteByCourseId(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程章节题目草稿失败")
	}
	if err := l.svcCtx.CourseCatalogueDraftModel.DeleteByCourseId(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程目录草稿失败")
	}
	if err := l.svcCtx.CourseTeacherDraftModel.DeleteByCourseId(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程老师草稿失败")
	}
	if err := l.svcCtx.CourseContentDraftModel.Delete(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程内容草稿失败")
	}
	if err := l.svcCtx.CourseDraftModel.Delete(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程草稿失败")
	}

	// ===== 正式表 =====
	catalogues, err := l.svcCtx.CourseCatalogueModel.ListByCourseId(l.ctx, id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	for _, c := range catalogues {
		if err := l.svcCtx.CourseCatalogueModel.Delete(l.ctx, c.Id); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程目录失败")
		}
	}
	if err := l.svcCtx.CourseSubjectModel.DeleteByCourseId(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程题目失败")
	}
	if err := l.svcCtx.CourseTeacherModel.DeleteByCourseId(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程老师失败")
	}
	if err := l.svcCtx.CourseContentModel.Delete(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程内容失败")
	}
	if err := l.svcCtx.CourseModel.Delete(l.ctx, id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除课程失败")
	}
	return &pb.Empty{}, nil
}
