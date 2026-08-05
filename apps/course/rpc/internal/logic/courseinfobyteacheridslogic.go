package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseInfoByTeacherIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseInfoByTeacherIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoByTeacherIdsLogic {
	return &CourseInfoByTeacherIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseInfoByTeacherIds 按老师 id 批量统计其关联的课程数量。
func (l *CourseInfoByTeacherIdsLogic) CourseInfoByTeacherIds(in *pb.TeacherIdsRequest) (*pb.TeacherCourseCountList, error) {
	if len(in.TeacherIds) == 0 {
		return &pb.TeacherCourseCountList{Items: []*pb.TeacherCourseCount{}}, nil
	}
	relations, err := l.svcCtx.CourseTeacherModel.ListByTeacherIds(l.ctx, in.TeacherIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询老师课程关系失败")
	}

	// 同一老师在同一课程下可能有多条记录，按课程 id 去重
	courseSet := make(map[int64]map[int64]struct{}, len(in.TeacherIds))
	for _, r := range relations {
		if _, ok := courseSet[r.TeacherId]; !ok {
			courseSet[r.TeacherId] = make(map[int64]struct{})
		}
		courseSet[r.TeacherId][r.CourseId] = struct{}{}
	}

	items := make([]*pb.TeacherCourseCount, 0, len(in.TeacherIds))
	for _, teacherId := range in.TeacherIds {
		items = append(items, &pb.TeacherCourseCount{
			TeacherId:  teacherId,
			CourseNum:  int64(len(courseSet[teacherId])),
			SubjectNum: 0, // 题目数由题库服务维护，course 库无对应数据
		})
	}
	return &pb.TeacherCourseCountList{Items: items}, nil
}
