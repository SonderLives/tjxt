// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeachersGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersGetLogic {
	return &CourseTeachersGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseTeachersGet 查询课程老师列表（透传 RPC）。
func (l *CourseTeachersGetLogic) CourseTeachersGet(req *types.CourseCataQueryReq) (resp []types.TeacherCourseInfoVO, err error) {
	list, err := l.svcCtx.CourseRpc.CourseTeachersGet(l.ctx, &pb.CourseTeachersGetRequest{
		Id:  req.Id,
		See: req.See,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程老师失败")
	}
	resp = make([]types.TeacherCourseInfoVO, 0, len(list.Items))
	for _, t := range list.Items {
		resp = append(resp, types.TeacherCourseInfoVO{
			Id:        t.Id,
			Name:      t.Name,
			Photo:     t.Photo,
			Job:       t.Job,
			Introduce: t.Introduce,
			Icon:      t.Icon,
			IsShow:    t.IsShow,
		})
	}
	return resp, nil
}
