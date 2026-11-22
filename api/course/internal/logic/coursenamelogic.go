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

type CourseNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseNameLogic {
	return &CourseNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseName 按课程名称模糊查询课程 id 列表（透传 RPC）。
func (l *CourseNameLogic) CourseName(req *types.CourseNameReq) (resp []int64, err error) {
	list, gerr := l.svcCtx.CourseRpc.CourseName2Ids(l.ctx, &pb.CourseNameRequest{Name: req.Name})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "按名称查询课程失败")
	}
	return list.Ids, nil
}
