package logic

import (
	"context"
	"database/sql"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSubjectsSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSubjectsSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsSaveLogic {
	return &CourseSubjectsSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseSubjectsSave 保存课程小节题目绑定：先清空旧绑定，再按 (course_id, cata_id, subject_id) 重新插入。
func (l *CourseSubjectsSaveLogic) CourseSubjectsSave(in *pb.CourseSubjectsSaveRequest) (*pb.Empty, error) {
	if in.CourseId == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	if err := l.svcCtx.CourseCataSubjectDraftModel.DeleteByCourseId(l.ctx, in.CourseId); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "清理课程题目关系失败")
	}

	courseId := sql.NullInt64{Int64: in.CourseId, Valid: true}
	for _, bind := range in.Subjects {
		if bind == nil {
			continue
		}
		for _, subjectId := range bind.SubjectIds {
			data := &model.CourseCataSubjectDraft{
				CourseId:  courseId,
				CataId:    bind.CataId,
				SubjectId: subjectId,
			}
			if _, err := l.svcCtx.CourseCataSubjectDraftModel.Insert(l.ctx, data); err != nil {
				return nil, xerr.Wrap(err, xerr.CodeInternal, "保存课程题目关系失败")
			}
		}
	}
	return &pb.Empty{}, nil
}
