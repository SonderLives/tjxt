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
// 老师姓名/头像/职位/简介来自 user 服务（GetUsersByIds 批量查询）；
// user 服务不可用时展示字段回落为空串，不阻塞课程主流程。
func (l *CourseTeachersGetLogic) CourseTeachersGet(in *pb.CourseTeachersGetRequest) (*pb.TeacherInfoList, error) {
	teachers, err := l.svcCtx.CourseTeacherModel.ListByCourseId(l.ctx, in.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}

	ids := make([]int64, 0, len(teachers))
	for _, t := range teachers {
		isShow := t.IsShow == 1
		// see 为 true 表示用户端查看，只返回允许展示的老师
		if in.See && !isShow {
			continue
		}
		ids = append(ids, t.TeacherId)
	}
	userMap := fetchTeacherMap(l.ctx, l.svcCtx, ids)

	items := make([]*pb.TeacherInfo, 0, len(teachers))
	for _, t := range teachers {
		isShow := t.IsShow == 1
		if in.See && !isShow {
			continue
		}
		ti := &pb.TeacherInfo{
			Id:     t.TeacherId,
			IsShow: isShow,
		}
		if u, ok := userMap[t.TeacherId]; ok && u != nil {
			ti.Name, ti.Photo, ti.Icon, ti.Job, ti.Introduce = u.Name, u.Photo, u.Icon, u.Job, u.Introduce
		}
		items = append(items, ti)
	}
	return &pb.TeacherInfoList{Items: items}, nil
}
