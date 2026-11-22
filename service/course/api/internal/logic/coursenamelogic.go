// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/service/course/rpc/pb"

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
		return nil, gerr
	}
	return list.Ids, nil
}
