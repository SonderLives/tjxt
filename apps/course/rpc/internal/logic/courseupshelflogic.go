package logic

import (
	"context"
	"database/sql"
	"time"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseUpShelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseUpShelfLogic {
	return &CourseUpShelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 上下架 =====
// CourseUpShelf 单课上架（发布）：把课程草稿及其目录/内容/老师复制到正式表，状态置为已上架。
func (l *CourseUpShelfLogic) CourseUpShelf(in *pb.IdRequest) (*pb.Empty, error) {
	if err := publishCourse(l.ctx, l.svcCtx, in.Id); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

// publishCourse 发布单个课程：course_draft -> course，并同步草稿子表到正式子表。
// 正式表主键沿用草稿主键（同一雪花 id）；草稿数据保留，支持再次编辑后重新发布。
// CourseUpShelf（单课）与 CourseUp（批量）共用。
func publishCourse(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	draft, err := svcCtx.CourseDraftModel.FindOne(ctx, id)
	if err != nil {
		if isNotFound(err) {
			return xerr.NotFound("课程草稿不存在")
		}
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}

	now := time.Now()
	course := &model.Course{
		Id:                draft.Id,
		Name:              draft.Name,
		CourseType:        draft.CourseType,
		CoverUrl:          draft.CoverUrl,
		FirstCateId:       draft.FirstCateId,
		SecondCateId:      draft.SecondCateId,
		ThirdCateId:       draft.ThirdCateId,
		Free:              draft.Free,
		Price:             draft.Price,
		TemplateType:      draft.TemplateType,
		TemplateUrl:       draft.TemplateUrl,
		Status:            int64(CourseStatusUpShelf),
		PurchaseStartTime: draft.PurchaseStartTime,
		PurchaseEndTime:   draft.PurchaseEndTime,
		Step:              draft.Step,
		Score:             draft.Score,
		MediaDuration:     sql.NullInt64{Int64: draft.MediaDuration, Valid: true},
		ValidDuration:     draft.ValidDuration,
		SectionNum:        sql.NullInt64{Int64: draft.SectionNum, Valid: true},
		DepId:             draft.DepId,
		PublishTimes:      1,
		PublishTime:       sql.NullTime{Time: now, Valid: true},
		Creater:           draft.Creater,
		Updater:           draft.Updater,
		Deleted:           draft.Deleted,
	}

	old, err := svcCtx.CourseModel.FindOne(ctx, id)
	switch {
	case err == nil:
		// 已发布过：累加发布次数，保留原始创建人
		course.PublishTimes = old.PublishTimes + 1
		course.Creater = old.Creater
		if err = svcCtx.CourseModel.Update(ctx, course); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "更新课程失败")
		}
	case isNotFound(err):
		if _, err = svcCtx.CourseModel.Insert(ctx, course); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "发布课程失败")
		}
	default:
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}

	if err = publishCourseContent(ctx, svcCtx, id, now); err != nil {
		return err
	}
	if err = publishCourseCatalogue(ctx, svcCtx, id); err != nil {
		return err
	}
	if err = publishCourseTeacher(ctx, svcCtx, id); err != nil {
		return err
	}
	// 发布上架事件，触发 search 服务增量同步 ES 课程索引（best-effort：失败仅告警）
	svcCtx.PublishCourseEvent(ctx, id, true)
	return nil
}

// publishCourseContent 课程内容草稿 -> 正式内容（按 id 替换插入；草稿缺失时写入空内容占位）。
func publishCourseContent(ctx context.Context, svcCtx *svc.ServiceContext, id int64, now time.Time) error {
	content := &model.CourseContent{
		Id:         id,
		CreateTime: now,
		UpdateTime: now,
	}
	draft, err := svcCtx.CourseContentDraftModel.FindOne(ctx, id)
	switch {
	case err == nil:
		content.CourseIntroduce = draft.CourseIntroduce
		content.UsePeople = draft.UsePeople
		content.CourseDetail = draft.CourseDetail
		content.DepId = draft.DepId
		content.Creater = draft.Creater
		content.Updater = draft.Updater
		content.Deleted = draft.Deleted
		if !draft.CreateTime.IsZero() {
			content.CreateTime = draft.CreateTime
		}
	case !isNotFound(err):
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程内容草稿失败")
	}
	if err = svcCtx.CourseContentModel.Upsert(ctx, content); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "保存课程内容失败")
	}
	return nil
}

// publishCourseCatalogue 课程目录草稿 -> 正式目录：先清空正式目录旧数据，再按草稿逐行复制（id 保持一致）。
func publishCourseCatalogue(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	olds, err := svcCtx.CourseCatalogueModel.ListByCourseId(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	for _, o := range olds {
		if err = svcCtx.CourseCatalogueModel.Delete(ctx, o.Id); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "清理旧课程目录失败")
		}
	}

	drafts, err := svcCtx.CourseCatalogueDraftModel.ListByCourseId(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录草稿失败")
	}
	for _, d := range drafts {
		data := &model.CourseCatalogue{
			Id:                d.Id,
			Name:              d.Name,
			Trailer:           d.Trailer,
			CourseId:          d.CourseId,
			Type:              d.Type,
			ParentCatalogueId: d.ParentCatalogueId,
			MediaId:           d.MediaId,
			VideoId:           d.VideoId,
			VideoName:         d.VideoName,
			LivingStartTime:   d.LivingStartTime,
			LivingEndTime:     d.LivingEndTime,
			PlayBack:          d.PlayBack,
			MediaDuration:     d.MediaDuration,
			CIndex:            d.CIndex,
			DepId:             d.DepId,
			Creater:           d.Creater,
			Updater:           d.Updater,
			Deleted:           d.Deleted,
		}
		if _, err = svcCtx.CourseCatalogueModel.Insert(ctx, data); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "发布课程目录失败")
		}
	}
	return nil
}

// publishCourseTeacher 课程老师草稿 -> 正式老师：先清空正式老师旧数据，再按草稿逐行复制（id 保持一致）。
func publishCourseTeacher(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	olds, err := svcCtx.CourseTeacherModel.ListByCourseId(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}
	for _, o := range olds {
		if err = svcCtx.CourseTeacherModel.Delete(ctx, o.Id); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "清理旧课程老师失败")
		}
	}

	drafts, err := svcCtx.CourseTeacherDraftModel.ListByCourseId(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程老师草稿失败")
	}
	for _, d := range drafts {
		data := &model.CourseTeacher{
			Id:        d.Id,
			CourseId:  d.CourseId,
			TeacherId: d.TeacherId,
			IsShow:    d.IsShow,
			CIndex:    d.CIndex,
			DepId:     d.DepId,
			Creater:   d.Creater,
			Updater:   d.Updater,
			Deleted:   d.Deleted,
		}
		if _, err = svcCtx.CourseTeacherModel.Insert(ctx, data); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "发布课程老师失败")
		}
	}
	return nil
}
