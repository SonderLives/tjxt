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

type CourseInfoByTeacherIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseInfoByTeacherIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoByTeacherIdsLogic {
	return &CourseInfoByTeacherIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseInfoByTeacherIds 按老师 id 统计课程/题目数量（透传 RPC）。
func (l *CourseInfoByTeacherIdsLogic) CourseInfoByTeacherIds(req *types.CourseInfoByTeacherIdsReq) (resp []types.TeacherCourseCountVO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CourseInfoByTeacherIds(l.ctx, &pb.TeacherIdsRequest{
		TeacherIds: parseIds(req.TeacherIds),
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询老师课程信息失败")
	}
	resp = make([]types.TeacherCourseCountVO, 0, len(list.Items))
	for _, item := range list.Items {
		resp = append(resp, types.TeacherCourseCountVO{
			TeacherId:  item.TeacherId,
			CourseNum:  item.CourseNum,
			SubjectNum: item.SubjectNum,
		})
	}
	return resp, nil
}
