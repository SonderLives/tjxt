package courselogic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"tjxt/pkg/errorx"
	"tjxt/app/course/rpc/internal/model"
	"tjxt/app/course/rpc/internal/svc"
	"tjxt/app/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
// 多表写入（course/course_content/course_catalogue/course_teacher）包进单个事务，
// 中途失败整体回滚，避免留下半发布状态；事务内只做本地 DB 写（无查询、无 RPC/MQ），
// 上架事件在事务提交后再发布。CourseUpShelf（单课）与 CourseUp（批量）共用。
func publishCourse(ctx context.Context, svcCtx *svc.ServiceContext, id int64) error {
	draft, err := svcCtx.CourseDraftModel.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return errorx.NotFound("课程草稿不存在")
		}
		return errorx.Wrap(err, errorx.CodeInternal, "查询课程草稿失败")
	}

	// 预读取全部数据，事务内只做纯写。
	old, err := svcCtx.CourseModel.FindOne(ctx, id)
	switch {
	case err == nil:
		// 已发布过：累加发布次数，保留原始创建人
	case errors.Is(err, model.ErrNotFound):
		old = nil
	default:
		return errorx.Wrap(err, errorx.CodeInternal, "查询课程失败")
	}
	course := buildPublishedCourse(draft, old, time.Now())

	content, err := buildPublishedContent(ctx, svcCtx, id, time.Now())
	if err != nil {
		return err
	}

	cataOldIds, catalogues, err := buildPublishedCatalogue(ctx, svcCtx, id)
	if err != nil {
		return err
	}

	teacherOldIds, teachers, err := buildPublishedTeachers(ctx, svcCtx, id)
	if err != nil {
		return err
	}

	if err = svcCtx.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if old == nil {
			if _, err := svcCtx.CourseModel.InsertTx(ctx, session, course); err != nil {
				return err
			}
		} else if err := svcCtx.CourseModel.UpdateTx(ctx, session, course); err != nil {
			return err
		}
		if err := svcCtx.CourseContentModel.UpsertTx(ctx, session, content); err != nil {
			return err
		}
		if err := svcCtx.CourseCatalogueModel.DeleteByIdsTx(ctx, session, cataOldIds); err != nil {
			return err
		}
		if err := svcCtx.CourseCatalogueModel.BatchInsertTx(ctx, session, catalogues); err != nil {
			return err
		}
		if err := svcCtx.CourseTeacherModel.DeleteByIdsTx(ctx, session, teacherOldIds); err != nil {
			return err
		}
		return svcCtx.CourseTeacherModel.BatchInsertTx(ctx, session, teachers)
	}); err != nil {
		return errorx.Wrap(err, errorx.CodeInternal, "发布课程写入失败")
	}

	// 发布上架事件，触发 search 服务增量同步 ES 课程索引（best-effort：失败仅告警）
	svcCtx.PublishCourseEvent(ctx, id, true)
	return nil
}

// buildPublishedCourse 依据课程草稿构造发布后的正式课程行（状态置为已上架）；
// old 非空表示已发布过：累加发布次数并保留原始创建人。
func buildPublishedCourse(draft *model.CourseDraft, old *model.Course, now time.Time) *model.Course {
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
	if old != nil {
		course.PublishTimes = old.PublishTimes + 1
		course.Creater = old.Creater
	}
	return course
}

// buildPublishedContent 课程内容草稿 -> 正式内容行（草稿缺失时写入空内容占位）。
func buildPublishedContent(ctx context.Context, svcCtx *svc.ServiceContext, id int64, now time.Time) (*model.CourseContent, error) {
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
	case !errors.Is(err, model.ErrNotFound):
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程内容草稿失败")
	}
	return content, nil
}

// buildPublishedCatalogue 课程目录草稿 -> 正式目录行：返回正式表旧行 id（待批量清理）
// 与草稿全量行（id 保持一致），删除/插入由事务内批量 SQL 完成。
func buildPublishedCatalogue(ctx context.Context, svcCtx *svc.ServiceContext, id int64) ([]int64, []*model.CourseCatalogue, error) {
	olds, err := svcCtx.CourseCatalogueModel.ListByCourseId(ctx, id)
	if err != nil {
		return nil, nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程目录失败")
	}
	oldIds := make([]int64, 0, len(olds))
	for _, o := range olds {
		oldIds = append(oldIds, o.Id)
	}

	drafts, err := svcCtx.CourseCatalogueDraftModel.ListByCourseId(ctx, id)
	if err != nil {
		return nil, nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程目录草稿失败")
	}
	list := make([]*model.CourseCatalogue, 0, len(drafts))
	for _, d := range drafts {
		list = append(list, &model.CourseCatalogue{
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
		})
	}
	return oldIds, list, nil
}

// buildPublishedTeachers 课程老师草稿 -> 正式老师行：返回正式表旧行 id（待批量清理）
// 与草稿全量行（id 保持一致），删除/插入由事务内批量 SQL 完成。
func buildPublishedTeachers(ctx context.Context, svcCtx *svc.ServiceContext, id int64) ([]int64, []*model.CourseTeacher, error) {
	olds, err := svcCtx.CourseTeacherModel.ListByCourseId(ctx, id)
	if err != nil {
		return nil, nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程老师失败")
	}
	oldIds := make([]int64, 0, len(olds))
	for _, o := range olds {
		oldIds = append(oldIds, o.Id)
	}

	drafts, err := svcCtx.CourseTeacherDraftModel.ListByCourseId(ctx, id)
	if err != nil {
		return nil, nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程老师草稿失败")
	}
	list := make([]*model.CourseTeacher, 0, len(drafts))
	for _, d := range drafts {
		list = append(list, &model.CourseTeacher{
			Id:        d.Id,
			CourseId:  d.CourseId,
			TeacherId: d.TeacherId,
			IsShow:    d.IsShow,
			CIndex:    d.CIndex,
			DepId:     d.DepId,
			Creater:   d.Creater,
			Updater:   d.Updater,
			Deleted:   d.Deleted,
		})
	}
	return oldIds, list, nil
}
