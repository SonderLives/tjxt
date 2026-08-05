package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// CourseStepCatalogue 课程目录已保存对应的填写进度
const CourseStepCatalogue = 2

type CourseCatalogueSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueSaveLogic {
	return &CourseCatalogueSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCatalogueSave 保存课程目录草稿：先清空该课程的旧草稿，再按树重新落库，最后推进草稿填写进度。
func (l *CourseCatalogueSaveLogic) CourseCatalogueSave(in *pb.CourseCatalogueSaveRequest) (*pb.Empty, error) {
	if in.CourseId == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	if err := l.svcCtx.CourseCatalogueDraftModel.DeleteByCourseId(l.ctx, in.CourseId); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "清理课程目录草稿失败")
	}
	if err := l.saveCataNodes(in.CourseId, 0, in.Chapters); err != nil {
		return nil, err
	}
	if err := l.advanceStep(in.CourseId, int64(in.Step)); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// saveCataNodes 递归写入目录草稿：章的 parent 为 0，其子节点（节/测试）的 parent 为所属章的 id。
// c_index 优先取请求中的 index，缺省时按同级顺序从 1 开始编号。
func (l *CourseCatalogueSaveLogic) saveCataNodes(courseId, parentId int64, nodes []*pb.CourseChapterSave) error {
	for i, n := range nodes {
		if n == nil {
			continue
		}
		cIndex := int64(i + 1)
		if n.Index > 0 {
			cIndex = int64(n.Index)
		}
		cataType := int64(n.Type)
		if cataType == 0 {
			if parentId == 0 {
				cataType = CatalogueTypeChapter
			} else {
				cataType = CatalogueTypeSection
			}
		}
		// 草稿表 id 为雪花 ID（非自增），需自行生成后再作为子节点的 parent
		id := nextID()
		data := &model.CourseCatalogueDraft{
			Id:                id,
			Name:              n.Name,
			CourseId:          courseId,
			Type:              cataType,
			ParentCatalogueId: parentId,
			CIndex:            cIndex,
			CanUpdate:         1,
		}
		if _, err := l.svcCtx.CourseCatalogueDraftModel.Insert(l.ctx, data); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "保存课程目录草稿失败")
		}
		if len(n.Sections) > 0 {
			if err := l.saveCataNodes(courseId, id, n.Sections); err != nil {
				return err
			}
		}
	}
	return nil
}

// advanceStep 推进草稿填写进度，目录保存后至少为 2，且不回退已有进度。
func (l *CourseCatalogueSaveLogic) advanceStep(courseId, step int64) error {
	if step < CourseStepCatalogue {
		step = CourseStepCatalogue
	}
	draft, err := l.svcCtx.CourseDraftModel.FindOne(l.ctx, courseId)
	if err != nil {
		if isNotFound(err) {
			return xerr.NotFound("课程草稿不存在")
		}
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	if draft.Step >= step {
		return nil
	}
	draft.Step = step
	if err := l.svcCtx.CourseDraftModel.Update(l.ctx, draft); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新课程草稿进度失败")
	}
	return nil
}
