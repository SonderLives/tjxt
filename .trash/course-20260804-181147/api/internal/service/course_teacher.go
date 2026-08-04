package service

import (
	"context"
	"database/sql"

	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"
	"tjxt/apps/course/api/internal/model"
	"tjxt/apps/course/api/internal/types"
)

// SaveTeachers 保存课程老师关系。
func (s *courseService) SaveTeachers(ctx context.Context, req *types.CourseTeacherSaveDTO, userId int64) error {
	// 1.课程草稿必须存在
	draft, err := s.courseDraftModel.FindById(ctx, req.Id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("课程草稿不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	// 2.组装并全删重插
	list := make([]*model.CourseTeacherDraft, 0, len(req.Teachers))
	for i, teacher := range req.Teachers {
		list = append(list, &model.CourseTeacherDraft{
			Id:        idgen.NextID(),
			CourseId:  req.Id,
			TeacherId: teacher.Id,
			IsShow:    boolToInt(teacher.IsShow),
			CIndex:    int64(i),
			DepId:     draft.DepId,
			Creater:   userId,
			Updater:   userId,
		})
	}
	if err := s.teacherDraftModel.ReplaceAll(ctx, req.Id, list); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "保存老师信息失败")
	}
	if err := s.UpdateStep(ctx, req.Id, CourseStepTeacher); err != nil {
		return err
	}
	return nil
}

// QueryTeachers 查询课程相关老师信息。
func (s *courseService) QueryTeachers(ctx context.Context, courseId int64, see bool) ([]*types.CourseTeacherVO, error) {
	if see {
		formal, err := s.queryFormalTeachers(ctx, courseId)
		if err != nil {
			return nil, err
		}
		if len(formal) > 0 {
			return formal, nil
		}
	}
	return s.queryDraftTeachers(ctx, courseId)
}

// queryFormalTeachers 查询正式课程老师信息。
func (s *courseService) queryFormalTeachers(ctx context.Context, courseId int64) ([]*types.CourseTeacherVO, error) {
	rels, err := s.teacherModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}
	if len(rels) == 0 {
		return []*types.CourseTeacherVO{}, nil
	}
	ids := make([]int64, 0, len(rels))
	for _, rel := range rels {
		ids = append(ids, rel.TeacherId)
	}
	return s.buildTeacherVOs(ctx, rels, ids)
}

// queryDraftTeachers 查询草稿课程老师信息。
func (s *courseService) queryDraftTeachers(ctx context.Context, courseId int64) ([]*types.CourseTeacherVO, error) {
	drafts, err := s.teacherDraftModel.ListByCourseId(ctx, courseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}
	if len(drafts) == 0 {
		return []*types.CourseTeacherVO{}, nil
	}
	rels := make([]*model.CourseTeacher, 0, len(drafts))
	ids := make([]int64, 0, len(drafts))
	for _, d := range drafts {
		rels = append(rels, &model.CourseTeacher{
			Id: d.Id, CourseId: d.CourseId, TeacherId: d.TeacherId,
			IsShow: d.IsShow, CIndex: d.CIndex, DepId: d.DepId,
			Creater: d.Creater, Updater: d.Updater,
		})
		ids = append(ids, d.TeacherId)
	}
	return s.buildTeacherVOs(ctx, rels, ids)
}

// buildTeacherVOs 组装老师信息 VO。
func (s *courseService) buildTeacherVOs(ctx context.Context, rels []*model.CourseTeacher, teacherIds []int64) ([]*types.CourseTeacherVO, error) {
	userMap, err := s.userDetailModel.FindByIds(ctx, distinctIDs(teacherIds))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询老师信息失败")
	}
	vos := make([]*types.CourseTeacherVO, 0, len(rels))
	for _, rel := range rels {
		vo := &types.CourseTeacherVO{Id: rel.TeacherId, IsShow: rel.IsShow == 1}
		if u, ok := userMap[rel.TeacherId]; ok {
			vo.Icon = u.Icon
			vo.Photo = u.Photo
			vo.Name = u.Name
			vo.Introduce = u.Intro
			vo.Job = u.Job
		}
		vos = append(vos, vo)
	}
	return vos, nil
}
