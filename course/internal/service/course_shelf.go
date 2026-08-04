package service

import (
	"context"
	"database/sql"
	"time"

	"common/xerr"
	"course/internal/model"
)

// UpdateStep 更新课程编辑进度，进度只能前进不能后退；目录/题目保存时刷新课时数。
func (s *courseService) UpdateStep(ctx context.Context, id int64, step int64) error {
	draft, err := s.courseDraftModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("课程草稿不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	update := &model.CourseDraft{
		Id:       id,
		CVersion: sql.NullInt64{Int64: draft.CVersion.Int64 + 1, Valid: true},
	}
	update.Step = draft.Step
	if draft.Step < step {
		update.Step = step
	}
	if step == CourseStepCatalogue || step == CourseStepSubject {
		num, err := s.cataDraftModel.CountSections(ctx, id)
		if err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "统计课时数失败")
		}
		update.SectionNum = num
	}
	if err := s.courseDraftModel.UpdateById(ctx, update, "step", "section_num", "c_version"); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新课程进度失败")
	}
	return nil
}

// CheckBeforeUpShelf 课程上架前校验。
func (s *courseService) CheckBeforeUpShelf(ctx context.Context, id int64) error {
	draft, err := s.courseDraftModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		_, e2 := s.courseModel.FindById(ctx, id)
		if e2 == nil {
			return xerr.Conflict("课程已上架")
		}
		return xerr.Conflict("课程不存在，无法上架")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	course, err := s.courseModel.FindById(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	// 草稿信息不完整无法上架
	if draft.Step != CourseStepTeacher {
		return xerr.Conflict("课程信息不完整，无法上架")
	}
	// 已上架/已完结课程不能上架
	if course != nil && course.Status != CourseStatusDownShelf {
		return xerr.Conflict("课程状态不允许上架")
	}
	// 校验课程结束时间
	if !draft.PurchaseEndTime.After(time.Now()) {
		return xerr.Conflict("课程结束时间已过，无法上架")
	}
	// 首次上架校验同名课程
	if course == nil {
		num, err := s.courseModel.CountSameName(ctx, draft.Name)
		if err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "校验课程名称失败")
		}
		if num > 0 {
			return xerr.BadRequestf("课程名称已存在，无法上架")
		}
	}
	// 校验课程目录（小节必须有媒资，练习必须有题目）
	if err := s.checkCataInfoComplete(ctx, id); err != nil {
		return err
	}
	return nil
}

// checkCataInfoComplete 校验目录完整性：小节已上传视频，练习已添加题目。
func (s *courseService) checkCataInfoComplete(ctx context.Context, courseId int64) error {
	catas, err := s.cataDraftModel.ListAllByCourseId(ctx, courseId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	if len(catas) == 0 {
		return xerr.Conflict("课程目录为空，无法上架")
	}
	subjects, err := s.cataSubjectDraftModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程题目关系失败")
	}
	subjectNumMap := make(map[int64]int64, len(subjects))
	for _, rel := range subjects {
		subjectNumMap[rel.CataId]++
	}
	for _, c := range catas {
		if c.Type == CataTypeSection && c.VideoName == "" {
			return xerr.Newf(xerr.CodeConflict, "小节《%s》未上传视频，无法上架", c.Name)
		}
		if c.Type == CataTypePractice && subjectNumMap[c.Id] <= 0 {
			return xerr.Newf(xerr.CodeConflict, "练习《%s》未添加题目，无法上架", c.Name)
		}
	}
	return nil
}

// UpShelf 课程上架。
func (s *courseService) UpShelf(ctx context.Context, id int64) error {
	// 1.校验
	if err := s.CheckBeforeUpShelf(ctx, id); err != nil {
		return err
	}
	draft, err := s.courseDraftModel.FindById(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	course, err := s.courseModel.FindById(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	isFirst := course == nil

	// 2.计算每个章节的课时视频时长
	mediaDurations, err := s.cataDraftModel.SumMediaDurationByCourse(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "统计课程视频时长失败")
	}
	totalMediaDuration := int64(0)
	for _, d := range mediaDurations {
		totalMediaDuration += d
	}

	// 3.目录与老师上架
	if err := s.copyTeacherToShelf(ctx, id); err != nil {
		return err
	}
	if err := s.copyCataToShelf(ctx, id); err != nil {
		return err
	}

	// 4.组装正式课程与内容
	contentDraft, err := s.courseContentDraftModel.FindById(ctx, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程内容草稿失败")
	}
	content := &model.CourseContent{
		Id:              id,
		CourseIntroduce: contentDraft.CourseIntroduce,
		UsePeople:       contentDraft.UsePeople,
		CourseDetail:    contentDraft.CourseDetail,
		DepId:           contentDraft.DepId,
		Creater:         contentDraft.Creater,
		Updater:         contentDraft.Updater,
	}
	now := time.Now()
	toShelf := &model.Course{
		Id:               id,
		Name:             draft.Name,
		CourseType:       draft.CourseType,
		CoverUrl:         draft.CoverUrl,
		FirstCateId:      draft.FirstCateId,
		SecondCateId:     draft.SecondCateId,
		ThirdCateId:      draft.ThirdCateId,
		Free:             draft.Free,
		Price:            draft.Price,
		TemplateType:     draft.TemplateType,
		TemplateUrl:      draft.TemplateUrl,
		Status:           CourseStatusShelf,
		PurchaseStartTime: draft.PurchaseStartTime,
		PurchaseEndTime:  draft.PurchaseEndTime,
		Step:             draft.Step,
		MediaDuration:    sql.NullInt64{Int64: totalMediaDuration, Valid: true},
		ValidDuration:    draft.ValidDuration,
		SectionNum:       sql.NullInt64{Int64: draft.SectionNum, Valid: true},
		DepId:            draft.DepId,
		PublishTime:      sql.NullTime{Time: now, Valid: true},
		Score:            sql.NullInt64{Int64: int64(40 + time.Now().UnixNano()%10), Valid: true},
		Creater:          draft.Creater,
		Updater:          draft.Updater,
	}
	if course != nil {
		toShelf.PublishTimes = sql.NullInt64{Int64: course.PublishTimes.Int64 + 1, Valid: true}
	} else {
		toShelf.PublishTimes = sql.NullInt64{Int64: 1, Valid: true}
	}

	if isFirst {
		if err := s.courseContentModel.Insert(ctx, content); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "新增课程内容失败")
		}
		if err := s.courseModel.Insert(ctx, toShelf); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "新增课程失败")
		}
	} else {
		if err := s.courseContentModel.UpdateById(ctx, content); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "更新课程内容失败")
		}
		if err := s.courseModel.UpdateById(ctx, toShelf,
			"name", "cover_url", "status", "media_duration", "valid_duration",
			"purchase_start_time", "purchase_end_time", "publish_time", "publish_times",
			"score", "section_num", "step"); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "更新课程失败")
		}
	}
	// 5.删除课程草稿与草稿内容
	if err := s.courseDraftModel.DeleteById(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程草稿失败")
	}
	if err := s.courseContentDraftModel.DeleteById(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程内容草稿失败")
	}
	return nil
}

// copyTeacherToShelf 将课程老师草稿上架到正式表。
func (s *courseService) copyTeacherToShelf(ctx context.Context, courseId int64) error {
	drafts, err := s.teacherDraftModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程老师草稿失败")
	}
	formal := make([]*model.CourseTeacher, 0, len(drafts))
	for _, d := range drafts {
		formal = append(formal, &model.CourseTeacher{
			Id:        d.Id,
			CourseId:  d.CourseId,
			TeacherId: d.TeacherId,
			IsShow:    d.IsShow,
			CIndex:    d.CIndex,
			DepId:     d.DepId,
			Creater:   d.Creater,
			Updater:   d.Updater,
		})
	}
	if err := s.teacherModel.SaveAll(ctx, courseId, formal); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "上架课程老师失败")
	}
	if err := s.teacherDraftModel.DeleteByCourseId(ctx, courseId); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程老师草稿失败")
	}
	return nil
}

// copyCataToShelf 将课程目录草稿上架到正式表。
func (s *courseService) copyCataToShelf(ctx context.Context, courseId int64) error {
	drafts, err := s.cataDraftModel.ListAllByCourseId(ctx, courseId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录草稿失败")
	}
	formal := make([]*model.CourseCatalogue, 0, len(drafts))
	for _, d := range drafts {
		formal = append(formal, &model.CourseCatalogue{
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
		})
	}
	if err := s.cataModel.SaveAll(ctx, formal); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "上架课程目录失败")
	}
	if err := s.cataDraftModel.DeleteByCourseId(ctx, courseId, []int64{CataTypeChapter, CataTypeSection, CataTypePractice}); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程目录草稿失败")
	}
	return nil
}

// DownShelf 课程下架（正式数据拷贝回草稿）。
func (s *courseService) DownShelf(ctx context.Context, id int64) error {
	course, err := s.courseModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return xerr.Conflict("课程不存在或未上架")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	if course.Status != CourseStatusShelf {
		return xerr.Conflict("课程状态不允许下架")
	}
	// 1.更新课程状态为下架
	if err := s.courseModel.UpdateById(ctx, &model.Course{Id: id, Status: CourseStatusDownShelf}, "status"); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新课程状态失败")
	}
	// 2.课程基本信息和内容信息拷贝到草稿
	if err := s.copyCourseToDraft(ctx, course); err != nil {
		return err
	}
	// 3.目录内容拷贝到草稿
	if err := s.copyCataToDraft(ctx, id); err != nil {
		return err
	}
	// 4.老师关系拷贝到草稿
	if err := s.copyTeacherToDraft(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *courseService) copyCourseToDraft(ctx context.Context, course *model.Course) error {
	content, err := s.courseContentModel.FindById(ctx, course.Id)
	if err != nil && err != sql.ErrNoRows {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程内容失败")
	}
	draft := &model.CourseDraft{
		Id:                course.Id,
		Name:              course.Name,
		CourseType:        course.CourseType,
		CoverUrl:          course.CoverUrl,
		FirstCateId:       course.FirstCateId,
		SecondCateId:      course.SecondCateId,
		ThirdCateId:       course.ThirdCateId,
		Free:              course.Free,
		Price:             course.Price,
		TemplateType:      course.TemplateType,
		TemplateUrl:       course.TemplateUrl,
		Status:            course.Status,
		PurchaseStartTime: course.PurchaseStartTime,
		PurchaseEndTime:   course.PurchaseEndTime,
		Step:              course.Step,
		Score:             course.Score,
		MediaDuration:     course.MediaDuration.Int64,
		ValidDuration:     course.ValidDuration,
		SectionNum:        course.SectionNum.Int64,
		CanUpdate:         1,
		DepId:             course.DepId,
		PublishTime:       course.PublishTime,
		Creater:           course.Creater,
		Updater:           course.Updater,
	}
	if err := s.courseDraftModel.Insert(ctx, draft); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "拷贝课程到草稿失败")
	}
	if content != nil {
		contentDraft := &model.CourseContentDraft{
			Id:              content.Id,
			CourseIntroduce: content.CourseIntroduce,
			UsePeople:       content.UsePeople,
			CourseDetail:    content.CourseDetail,
			DepId:           content.DepId,
			Creater:         content.Creater,
			Updater:         content.Updater,
		}
		if err := s.courseContentDraftModel.Insert(ctx, contentDraft); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "拷贝课程内容到草稿失败")
		}
	}
	return nil
}

func (s *courseService) copyCataToDraft(ctx context.Context, courseId int64) error {
	formal, err := s.cataModel.ListByCourseId(ctx, courseId, true)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	drafts := make([]*model.CourseCatalogueDraft, 0, len(formal))
	for _, c := range formal {
		drafts = append(drafts, &model.CourseCatalogueDraft{
			Id:                c.Id,
			Name:              c.Name,
			Trailer:           c.Trailer,
			CourseId:          c.CourseId,
			Type:              c.Type,
			ParentCatalogueId: c.ParentCatalogueId,
			MediaId:           c.MediaId,
			VideoId:           c.VideoId,
			VideoName:         c.VideoName,
			LivingStartTime:   c.LivingStartTime,
			LivingEndTime:     c.LivingEndTime,
			PlayBack:          c.PlayBack,
			CIndex:            c.CIndex,
			MediaDuration:     c.MediaDuration,
			CanUpdate:         1,
			DepId:             c.DepId,
			Creater:           c.Creater,
			Updater:           c.Updater,
		})
	}
	if err := s.cataDraftModel.ReplaceAll(ctx, courseId, drafts); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "拷贝目录到草稿失败")
	}
	return nil
}

func (s *courseService) copyTeacherToDraft(ctx context.Context, courseId int64) error {
	formal, err := s.teacherModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}
	drafts := make([]*model.CourseTeacherDraft, 0, len(formal))
	for _, c := range formal {
		drafts = append(drafts, &model.CourseTeacherDraft{
			Id:        c.Id,
			CourseId:  c.CourseId,
			TeacherId: c.TeacherId,
			IsShow:    c.IsShow,
			CIndex:    c.CIndex,
			DepId:     c.DepId,
			Creater:   c.Creater,
			Updater:   c.Updater,
		})
	}
	if err := s.teacherDraftModel.ReplaceAll(ctx, courseId, drafts); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "拷贝老师到草稿失败")
	}
	return nil
}
