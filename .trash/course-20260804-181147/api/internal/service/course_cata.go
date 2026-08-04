package service

import (
	"context"
	"database/sql"
	"sort"

	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"
	"tjxt/apps/course/api/internal/model"
	"tjxt/apps/course/api/internal/types"
)

// SaveCatas 保存课程章节目录。
func (s *courseService) SaveCatas(ctx context.Context, courseId int64, step int64, list []*types.CataSaveDTO) error {
	// 1.按章序号升序排序
	sort.SliceStable(list, func(i, j int) bool { return list[i].Index < list[j].Index })
	// 2.校验章序号
	if len(list) == 0 {
		return xerr.BadRequestf("目录不能为空")
	}
	seen := make(map[int32]bool, len(list))
	for _, c := range list {
		if seen[c.Index] {
			return xerr.Conflict("章序号存在重复")
		}
		seen[c.Index] = true
	}
	if int(list[len(list)-1].Index) > len(list) {
		return xerr.Conflict("章序号存在间断")
	}
	if step != CourseStepCatalogue && step != CourseStepSubject {
		return xerr.BadRequestf("非法操作")
	}
	// 3.查询已上架目录
	shelfCatas, err := s.cataModel.ListByCourseId(ctx, courseId, true)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	if err := checkCataIndex(list, shelfCatas); err != nil {
		return err
	}
	// 4.组装数据（优先级：本次保存 > 草稿 > 已上架）
	drafts, err := s.packageCatas(ctx, courseId, list, shelfCatas)
	if err != nil {
		return err
	}
	// 5.删除原目录并重新插入
	typesToDelete := []int64{CataTypeChapter, CataTypeSection}
	if step == CourseStepSubject {
		typesToDelete = []int64{CataTypeChapter, CataTypeSection, CataTypePractice}
	}
	if err := s.cataDraftModel.DeleteByCourseId(ctx, courseId, typesToDelete); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程目录失败")
	}
	if err := s.cataDraftModel.ReplaceAll(ctx, courseId, drafts); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "保存课程目录失败")
	}
	// 6.更新编辑进度
	if err := s.UpdateStep(ctx, courseId, CourseStepCatalogue); err != nil {
		return err
	}
	// 7.删除已删章节的题目关系
	cataIds := make([]int64, 0, len(drafts))
	for _, d := range drafts {
		cataIds = append(cataIds, d.Id)
	}
	if err := s.cataSubjectDraftModel.DeleteNotInCataIds(ctx, courseId, cataIds); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "清理题目关系失败")
	}
	return nil
}

// checkCataIndex 校验已上架章目录是否被移动或删除。
func checkCataIndex(list []*types.CataSaveDTO, shelfCatas []*model.CourseCatalogue) error {
	// 章 id -> 新序号
	saveIndexMap := make(map[int64]int64, len(list))
	for _, c := range list {
		if c.Id > 0 {
			saveIndexMap[c.Id] = int64(c.Index)
		}
	}
	for _, sc := range shelfCatas {
		if sc.Type != CataTypeChapter {
			continue
		}
		index, ok := saveIndexMap[sc.Id]
		if !ok {
			return xerr.Newf(xerr.CodeConflict, "已上架的章《%s》已被删除，无法保存", sc.Name)
		}
		if index != sc.Index {
			return xerr.Newf(xerr.CodeConflict, "已上架的章《%s》序号被修改，无法保存", sc.Name)
		}
	}
	return nil
}

// packageCatas 组装目录草稿数据：本次保存 > 草稿 > 已上架。
func (s *courseService) packageCatas(ctx context.Context, courseId int64, list []*types.CataSaveDTO, shelfCatas []*model.CourseCatalogue) ([]*model.CourseCatalogueDraft, error) {
	savedMap := make(map[int64]*model.CourseCatalogueDraft)
	for _, sc := range shelfCatas {
		savedMap[sc.Id] = &model.CourseCatalogueDraft{
			Id:            sc.Id,
			Name:          sc.Name,
			Trailer:       sc.Trailer,
			CourseId:      sc.CourseId,
			Type:          sc.Type,
			ParentId:      sc.ParentId,
			MediaId:       sc.MediaId,
			MediaName:     sc.MediaName,
			Index:         sc.Index,
			MediaDuration: sc.MediaDuration,
			CanUpdate:     1,
		}
	}
	draftCatas, err := s.cataDraftModel.ListAllByCourseId(ctx, courseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录草稿失败")
	}
	for _, dc := range draftCatas {
		savedMap[dc.Id] = dc
	}

	result := make([]*model.CourseCatalogueDraft, 0, len(list))
	for _, cata := range list {
		chapterId := cata.Id
		if chapterId == 0 {
			chapterId = idgen.NextID()
		}
chapter := savedMap[chapterId]
	if chapter == nil {
		chapter = &model.CourseCatalogueDraft{Id: chapterId}
	}
	chapter.Name = cata.Name
	chapter.Type = CataTypeChapter
	chapter.CourseId = courseId
	chapter.ParentId = 0
	chapter.Index = int64(cata.Index)
	result = append(result, chapter)

	// 小节序号从 1 开始，练习无序号
	sectionIndex := int64(0)
	for _, section := range cata.Sections {
		sectionId := section.Id
		if sectionId == 0 {
			sectionId = idgen.NextID()
		}
		item := savedMap[sectionId]
		if item == nil {
			item = &model.CourseCatalogueDraft{Id: sectionId}
		}
		item.Name = section.Name
		item.Type = int64(section.Type)
		item.CourseId = courseId
		item.ParentId = chapterId
		if section.Type == int32(CataTypeSection) {
			sectionIndex++
			item.Index = sectionIndex
		}
		result = append(result, item)
	}
	}
	return result, nil
}

// QueryCatas 获取课程章节（see=true 优先正式，否则草稿）。
func (s *courseService) QueryCatas(ctx context.Context, courseId int64, see, withPractice bool) ([]*types.CataVO, error) {
	if see {
		formal, err := s.queryFormalCatasVO(ctx, courseId, withPractice)
		if err != nil {
			return nil, err
		}
		if len(formal) > 0 {
			return formal, nil
		}
	}
	return s.queryDraftCatasVO(ctx, courseId, withPractice)
}

// queryDraftCatasVO 查询并组装草稿目录树。
func (s *courseService) queryDraftCatasVO(ctx context.Context, courseId int64, withPractice bool) ([]*types.CataVO, error) {
	catas, err := s.cataDraftModel.ListByCourseId(ctx, courseId, withPractice)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	if len(catas) == 0 {
		return []*types.CataVO{}, nil
	}
	draft, err := s.courseDraftModel.FindById(ctx, courseId)
	if err != nil && err != sql.ErrNoRows {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	isNoUpShelf := draft != nil && draft.Status == CourseStatusNoUpShelf

	// 已上架部分的最大章序号、每章最大小节序号
	maxChapterIndex := int64(0)
	chapterIdAndMaxSection := make(map[int64]int64)
	if !isNoUpShelf {
		for _, c := range catas {
			if c.Type == CataTypeSection && c.CanUpdate == 0 {
				if c.Index > chapterIdAndMaxSection[c.ParentId] {
					chapterIdAndMaxSection[c.ParentId] = c.Index
				}
			}
		}
		for _, c := range catas {
			if c.Type == CataTypeChapter && c.CanUpdate == 0 && c.Index > maxChapterIndex {
				maxChapterIndex = c.Index
			}
		}
	}

	// 题目数量
	subjects, err := s.cataSubjectDraftModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程题目关系失败")
	}
	subjectNumMap := make(map[int64]int64, len(subjects))
	for _, rel := range subjects {
		subjectNumMap[rel.CataId]++
	}

	flats := make([]*cataFlat, 0, len(catas))
	for _, c := range catas {
		flats = append(flats, &cataFlat{
			Id:                c.Id,
			ParentId:          c.ParentId,
			Type:              c.Type,
			CIndex:            c.Index,
			Name:              c.Name,
			Trailer:           c.Trailer == 1,
			MediaDuration:     c.MediaDuration,
			VideoName:         c.MediaName,
			MediaId:           c.MediaId,
			CanUpdate:         c.CanUpdate == 1,
			SubjectNum:        subjectNumMap[c.Id],
			MaxIndexOnShelf:   computeMaxIndexOnShelf(c, chapterIdAndMaxSection, maxChapterIndex),
			MaxSectionIndexOnShelf: computeMaxSectionIndexOnShelf(c, chapterIdAndMaxSection, maxChapterIndex),
		})
	}
	return buildCataTree(flats), nil
}

func computeMaxIndexOnShelf(c *model.CourseCatalogueDraft, maxSection map[int64]int64, maxChapter int64) int64 {
	if c.Type == CataTypeSection {
		return maxSection[c.ParentId]
	}
	if c.Type == CataTypeChapter {
		return maxChapter
	}
	return 0
}

func computeMaxSectionIndexOnShelf(c *model.CourseCatalogueDraft, maxSection map[int64]int64, maxChapter int64) int64 {
	if c.Type == CataTypeChapter {
		return maxSection[c.Id]
	}
	return 0
}

// queryFormalCatasVO 查询并组装正式目录树。
func (s *courseService) queryFormalCatasVO(ctx context.Context, courseId int64, withPractice bool) ([]*types.CataVO, error) {
	catas, err := s.cataModel.ListByCourseId(ctx, courseId, withPractice)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	if len(catas) == 0 {
		return []*types.CataVO{}, nil
	}
	flats := make([]*cataFlat, 0, len(catas))
	for _, c := range catas {
		flats = append(flats, &cataFlat{
			Id:            c.Id,
			ParentId:      c.ParentId,
			Type:          c.Type,
			CIndex:        c.Index,
			Name:          c.Name,
			Trailer:       c.Trailer == 1,
			MediaDuration: c.MediaDuration,
			VideoName:     c.MediaName,
			MediaId:       c.MediaId,
			CanUpdate:     false,
		})
	}
	return buildCataTree(flats), nil
}

// cataFlat 目录树扁平节点。
type cataFlat struct {
	Id                    int64
	ParentId              int64
	Type                  int64
	CIndex                int64
	Name                  string
	Trailer               bool
	MediaDuration         int64
	VideoName             string
	MediaId               int64
	CanUpdate             bool
	SubjectNum            int64
	TotalScore            int64
	MaxIndexOnShelf       int64
	MaxSectionIndexOnShelf int64
}

// buildCataTree 将扁平目录列表组装为树。
func buildCataTree(list []*cataFlat) []*types.CataVO {
	nodeMap := make(map[int64]*types.CataVO, len(list))
	for _, f := range list {
		nodeMap[f.Id] = &types.CataVO{
			Id:                   f.Id,
			Index:                int32(f.CIndex),
			Name:                 f.Name,
			MediaDuration:        int32(f.MediaDuration),
			Trailer:              f.Trailer,
			MediaName:            f.VideoName,
			MediaId:              f.MediaId,
			Type:                 int32(f.Type),
			SubjectNum:           int32(f.SubjectNum),
			TotalScore:           int32(f.TotalScore),
			CanUpdate:            f.CanUpdate,
			MaxIndexOnShelf:      int32(f.MaxIndexOnShelf),
			MaxSectionIndexOnShelf: int32(f.MaxSectionIndexOnShelf),
		}
	}
	roots := make([]*types.CataVO, 0)
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

// SaveMedia 保存小节媒资信息。
func (s *courseService) SaveMedia(ctx context.Context, courseId int64, list []*types.CourseMediaSaveDTO) error {
	// 1.校验小节 id 属于当前课程
	cataIds := make([]int64, 0, len(list))
	for _, m := range list {
		cataIds = append(cataIds, m.CataId)
	}
	if err := s.checkSectionIds(ctx, cataIds, courseId); err != nil {
		return err
	}
	// 2.必须先保存目录
	draft, err := s.courseDraftModel.FindById(ctx, courseId)
	if err == sql.ErrNoRows || draft == nil || draft.Step < CourseStepCatalogue {
		return xerr.Conflict("请先保存课程目录")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	// 3.更新小节媒资信息
	items := make([]*model.CourseCatalogueDraft, 0, len(list))
	for _, m := range list {
		items = append(items, &model.CourseCatalogueDraft{
			Id:            m.CataId,
			MediaId:       m.MediaId,
			Trailer:       boolToInt(m.Trailer),
			MediaName:     m.VideoName,
			MediaDuration: int64(m.MediaDuration),
		})
	}
	if err := s.cataDraftModel.SaveMediaInfo(ctx, items); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "保存小节媒资失败")
	}
	if err := s.UpdateStep(ctx, courseId, CourseStepMedia); err != nil {
		return err
	}
	// 4.统计并更新每个章的媒资总时长
	durations, err := s.cataDraftModel.SumMediaDurationByChapter(ctx, courseId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "统计章节时长失败")
	}
	if err := s.cataDraftModel.UpdateMediaDuration(ctx, durations); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新章节时长失败")
	}
	return nil
}

// checkSectionIds 校验小节 id 列表与数据库小节集合一致。
func (s *courseService) checkSectionIds(ctx context.Context, cataIds []int64, courseId int64) error {
	sections, err := s.cataDraftModel.ListByTypes(ctx, courseId, []int64{CataTypeSection})
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程小节失败")
	}
	if len(sections) != len(cataIds) {
		return xerr.Conflict("所有小节都必须上传视频")
	}
	dbIds := make(map[int64]bool, len(sections))
	for _, sec := range sections {
		dbIds[sec.Id] = true
	}
	for _, id := range cataIds {
		if !dbIds[id] {
			return xerr.Conflict("小节数据非法")
		}
	}
	return nil
}

// SaveSubjects 保存小节/练习的题目。
func (s *courseService) SaveSubjects(ctx context.Context, courseId int64, list []*types.CataSubjectDTO) error {
	// 1.校验目录 id 属于当前课程的小节/练习
	cataIds := make([]int64, 0, len(list))
	for _, c := range list {
		cataIds = append(cataIds, c.CataId)
	}
	if err := s.checkPracticeIds(ctx, cataIds, courseId); err != nil {
		return err
	}
	// 2.必须先上传媒资
	draft, err := s.courseDraftModel.FindById(ctx, courseId)
	if err == sql.ErrNoRows || draft == nil || draft.Step < CourseStepMedia {
		return xerr.Conflict("请先上传课程视频")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	// 3.删除原有关系并插入
	if err := s.cataSubjectDraftModel.DeleteByCourseId(ctx, courseId); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除题目关系失败")
	}
	for _, c := range list {
		for _, subjectId := range c.SubjectIds {
			if _, err := s.cataSubjectDraftModel.Insert(ctx, &model.CourseCataSubjectDraft{
				Id:        idgen.NextID(),
				CourseId:  sql.NullInt64{Int64: courseId, Valid: true},
				CataId:    c.CataId,
				SubjectId: subjectId,
			}); err != nil {
				return xerr.Wrap(err, xerr.CodeInternal, "保存题目关系失败")
			}
		}
	}
	if err := s.UpdateStep(ctx, courseId, CourseStepSubject); err != nil {
		return err
	}
	return nil
}

// checkPracticeIds 校验目录 id 列表是课程小节/练习集合的子集，且练习都传入。
func (s *courseService) checkPracticeIds(ctx context.Context, cataIds []int64, courseId int64) error {
	items, err := s.cataDraftModel.ListByTypes(ctx, courseId, []int64{CataTypeSection, CataTypePractice})
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程目录失败")
	}
	if len(items) == 0 {
		return xerr.Conflict("课程目录不存在")
	}
	allIds := make(map[int64]bool, len(items))
	practiceIds := make([]int64, 0)
	for _, item := range items {
		allIds[item.Id] = true
		if item.Type == CataTypePractice {
			practiceIds = append(practiceIds, item.Id)
		}
	}
	submitted := make(map[int64]bool, len(cataIds))
	for _, id := range cataIds {
		if !allIds[id] {
			return xerr.Conflict("传入的目录不属于当前课程")
		}
		submitted[id] = true
	}
	for _, pid := range practiceIds {
		if !submitted[pid] {
			return xerr.Conflict("所有练习都必须添加题目")
		}
	}
	return nil
}

// GetSubjects 获取小节/练习中的题目。
func (s *courseService) GetSubjects(ctx context.Context, courseId int64) ([]*types.CataSimpleSubjectVO, error) {
	rels, err := s.cataSubjectDraftModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询题目关系失败")
	}
	if len(rels) == 0 {
		return []*types.CataSimpleSubjectVO{}, nil
	}
	subjectIds := make([]int64, 0, len(rels))
	for _, rel := range rels {
		subjectIds = append(subjectIds, rel.SubjectId)
	}
	subjects, err := s.subjectModel.ListByIds(ctx, distinctIDs(subjectIds))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询题目失败")
	}
	nameMap := make(map[int64]string, len(subjects))
	for _, sub := range subjects {
		nameMap[sub.Id] = sub.Content
	}
	// 按目录分组
	group := make(map[int64][]*types.SubjectInfo)
	for _, rel := range rels {
		group[rel.CataId] = append(group[rel.CataId], &types.SubjectInfo{
			Id:   rel.SubjectId,
			Name: nameMap[rel.SubjectId],
		})
	}
	result := make([]*types.CataSimpleSubjectVO, 0, len(group))
	for cataId, infos := range group {
		result = append(result, &types.CataSimpleSubjectVO{CataId: cataId, Subjects: infos})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CataId < result[j].CataId })
	return result, nil
}
