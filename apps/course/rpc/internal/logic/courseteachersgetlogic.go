package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseTeachersGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersGetLogic {
	return &CourseTeachersGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程老师 =====
// CourseTeachersGet 查询课程关联的老师。
// 注意：老师姓名/头像/职位/简介来自 user 服务，course 服务当前未接入 user 客户端，
// 故这里只返回 course_teacher 自身字段（老师 id 与是否展示），其余字段留空由上层补齐。
func (l *CourseTeachersGetLogic) CourseTeachersGet(in *pb.CourseTeachersGetRequest) (*pb.TeacherInfoList, error) {
	teachers, err := l.svcCtx.CourseTeacherModel.ListByCourseId(l.ctx, in.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}

	items := make([]*pb.TeacherInfo, 0, len(teachers))
	for _, t := range teachers {
		isShow := t.IsShow == 1
		// see 为 true 表示用户端查看，只返回允许展示的老师
		if in.See && !isShow {
			continue
		}
		items = append(items, &pb.TeacherInfo{
			Id:     t.TeacherId,
			IsShow: isShow,
		})
	}
	return &pb.TeacherInfoList{Items: items}, nil
}
