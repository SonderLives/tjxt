package service

import (
	"context"
	"database/sql"

	"common/xerr"
	"course/internal/types"
)

// GetInfoById 获取课程详细信息（含目录、老师）。
func (s *courseService) GetInfoById(ctx context.Context, id int64, withCatalogue, withTeachers bool) (*types.CourseFullInfoDTO, error) {
	course, err := s.courseModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("课程不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	vo := &types.CourseFullInfoDTO{
		Id:              course.Id,
		Name:            course.Name,
		Price:           int32(course.Price),
		ValidDuration:   int32(course.ValidDuration),
		CoverUrl:        course.CoverUrl,
		PurchaseEndTime: course.PurchaseEndTime.Format(timeFormat),
		FirstCateId:     course.FirstCateId,
		SecondCateId:    course.SecondCateId,
		ThirdCateId:     course.ThirdCateId,
		SectionNum:      int32(course.SectionNum.Int64),
	}
	if withCatalogue {
		chapters, err := s.queryFormalCatalogueDTO(ctx, id, true)
		if err != nil {
			return nil, err
		}
		vo.Chapters = chapters
	}
	if withTeachers {
		rels, err := s.teacherModel.ListByCourseId(ctx, id)
		if err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
		}
		teacherIds := make([]int64, 0, len(rels))
		for _, rel := range rels {
			teacherIds = append(teacherIds, rel.TeacherId)
		}
		vo.TeacherIds = teacherIds
	}
	return vo, nil
}

// queryFormalCatalogueDTO 组装正式目录树（内部 CourseFullInfoDTO 使用）。
func (s *courseService) queryFormalCatalogueDTO(ctx context.Context, courseId int64, withPractice bool) ([]*types.CatalogueDTO, error) {
	catas, err := s.cataModel.ListByCourseId(ctx, courseId, withPractice)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	if len(catas) == 0 {
		return []*types.CatalogueDTO{}, nil
	}
	flats := make([]*cataFlat, 0, len(catas))
	for _, c := range catas {
		flats = append(flats, &cataFlat{
			Id:            c.Id,
			ParentId:      c.ParentCatalogueId,
			Type:          c.Type,
			CIndex:        c.CIndex,
			Name:          c.Name,
			Trailer:       c.Trailer == 1,
			MediaDuration: c.MediaDuration,
			VideoName:     c.VideoName,
			MediaId:       c.MediaId,
		})
	}
	return buildCatalogueTree(flats), nil
}

// buildCatalogueTree 将扁平目录列表组装为内部 CatalogueDTO 树。
func buildCatalogueTree(list []*cataFlat) []*types.CatalogueDTO {
	nodeMap := make(map[int64]*types.CatalogueDTO, len(list))
	for _, f := range list {
		nodeMap[f.Id] = &types.CatalogueDTO{
			Id:            f.Id,
			Index:         int32(f.CIndex),
			Name:          f.Name,
			MediaDuration: int32(f.MediaDuration),
			Trailer:       f.Trailer,
			MediaName:     f.VideoName,
			MediaId:       f.MediaId,
			Type:          int32(f.Type),
		}
	}
	roots := make([]*types.CatalogueDTO, 0)
	for _, f := range list {
		vo := nodeMap[f.Id]
		if f.ParentId == 0 {
			roots = append(roots, vo)
		} else if parent, ok := nodeMap[f.ParentId]; ok {
			parent.Sections = append(parent.Sections, vo)
		} else {
			roots = append(roots, vo)
		}
	}
	return roots
}

// GetCourseDTOById 获取课程信息（搜索索引库使用）。
func (s *courseService) GetCourseDTOById(ctx context.Context, id int64) (*types.CourseDTO, error) {
	// 先查询正式课程，未找到再查草稿
	course, err := s.courseModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return s.courseDTOFromDraft(ctx, id)
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	// 查询课程第一位老师
	teacherId := int64(0)
	rels, err := s.teacherModel.ListByCourseId(ctx, id)
	if err == nil && len(rels) > 0 {
		teacherId = rels[0].TeacherId
	}
	status := ""
	switch course.Status {
	case CourseStatusNoUpShelf:
		status = "待上架"
	case CourseStatusShelf:
		status = "已上架"
	case CourseStatusDownShelf:
		status = "下架"
	case CourseStatusFinished:
		status = "已完结"
	}
	vo := &types.CourseDTO{
		Id:            course.Id,
		Name:          course.Name,
		CoverUrl:      course.CoverUrl,
		Price:         int32(course.Price),
		Free:          intToBool(course.Free),
		Status:        status,
		Step:          int32(course.Step),
		Score:         int32(course.Score.Int64),
		Sections:      int32(course.SectionNum.Int64),
		Duration:      int32(course.MediaDuration.Int64),
		PublishTime:   formatNullTime(course.PublishTime),
		CreateTime:    course.CreateTime.Format(timeFormat),
		UpdateTime:    course.UpdateTime.Format(timeFormat),
		Updater:       course.Updater,
		Teacher:       teacherId,
		CourseType:    int32(course.CourseType),
		Enable:        int32(course.Status),
		ValidDuration: int32(course.ValidDuration),
		CategoryIdLv1: course.FirstCateId,
		CategoryIdLv2: course.SecondCateId,
		CategoryIdLv3: course.ThirdCateId,
	}
	return vo, nil
}

// courseDTOFromDraft 从草稿组装 CourseDTO。
func (s *courseService) courseDTOFromDraft(ctx context.Context, id int64) (*types.CourseDTO, error) {
	draft, err := s.courseDraftModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return &types.CourseDTO{}, nil
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	teacherId := int64(0)
	rels, err := s.teacherDraftModel.ListByCourseId(ctx, id)
	if err == nil && len(rels) > 0 {
		teacherId = rels[0].TeacherId
	}
	vo := &types.CourseDTO{
		Id:            draft.Id,
		Name:          draft.Name,
		CoverUrl:      draft.CoverUrl,
		Price:         int32(draft.Price),
		Free:          intToBool(draft.Free),
		Status:        "待上架",
		Step:          int32(draft.Step),
		Sections:      int32(draft.SectionNum),
		Duration:      int32(draft.MediaDuration),
		CreateTime:    draft.CreateTime.Format(timeFormat),
		UpdateTime:    draft.UpdateTime.Format(timeFormat),
		Updater:       draft.Updater,
		Teacher:       teacherId,
		CourseType:    int32(draft.CourseType),
		Enable:        int32(draft.Status),
		ValidDuration: int32(draft.ValidDuration),
		CategoryIdLv1: draft.FirstCateId,
		CategoryIdLv2: draft.SecondCateId,
		CategoryIdLv3: draft.ThirdCateId,
	}
	return vo, nil
}

// QueryCourseAndCatalog 查询课程基本信息、目录、学习进度。
func (s *courseService) QueryCourseAndCatalog(ctx context.Context, id int64) (*types.CourseAndSectionVO, error) {
	full, err := s.GetInfoById(ctx, id, true, true)
	if err != nil {
		return nil, err
	}
	if full == nil {
		return nil, xerr.NotFound("课程不存在")
	}
	vo := &types.CourseAndSectionVO{
		Id:       id,
		Name:     full.Name,
		Sections: full.SectionNum,
		CoverUrl: full.CoverUrl,
	}
	// 老师信息
	if len(full.TeacherIds) > 0 {
		userMap, err := s.userDetailModel.FindByIds(ctx, full.TeacherIds)
		if err == nil {
			for _, tid := range full.TeacherIds {
				if u, ok := userMap[tid]; ok {
					vo.TeacherName = u.Name
					vo.TeacherIcon = u.Icon
					break
				}
			}
		}
	}
	// 组装章节
	chapters := make([]*types.ChapterVO, 0, len(full.Chapters))
	for _, c := range full.Chapters {
		cv := &types.ChapterVO{
			Id:            c.Id,
			Index:         c.Index,
			Name:          c.Name,
			MediaDuration: c.MediaDuration,
			Sections:      make([]*types.SectionVO, 0, len(c.Sections)),
		}
		for _, sec := range c.Sections {
			cv.Sections = append(cv.Sections, &types.SectionVO{
				Id:            sec.Id,
				Name:          sec.Name,
				Index:         sec.Index,
				Type:          sec.Type,
				MediaDuration: sec.MediaDuration,
				MediaId:       sec.MediaId,
				Trailer:       sec.Trailer,
				SubjectNum:    sec.SubjectNum,
				HasTest:       sec.SubjectNum > 0,
			})
		}
		chapters = append(chapters, cv)
	}
	vo.Chapters = chapters
	// 学习进度：learning 服务后续接入，当前返回默认值
	return vo, nil
}

// CountSubjectNumAndCourseNumOfTeacher 统计老师名下课程数与出题数。
func (s *courseService) CountSubjectNumAndCourseNumOfTeacher(ctx context.Context, teacherIds []int64) ([]*types.SubNumAndCourseNumDTO, error) {
	if len(teacherIds) == 0 {
		return []*types.SubNumAndCourseNumDTO{}, nil
	}
	formalMap, err := s.teacherModel.CountTeacherCourse(ctx, teacherIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计老师课程失败")
	}
	draftMap, err := s.teacherDraftModel.CountTeacherCourse(ctx, teacherIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计老师课程失败")
	}
	result := make([]*types.SubNumAndCourseNumDTO, 0, len(teacherIds))
	for _, tid := range teacherIds {
		result = append(result, &types.SubNumAndCourseNumDTO{
			TeacherId:  tid,
			CourseNum:  int32(formalMap[tid] + draftMap[tid]),
			SubjectNum: 0, // exam 服务后续接入
		})
	}
	return result, nil
}
