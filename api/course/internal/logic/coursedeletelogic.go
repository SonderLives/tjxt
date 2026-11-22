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

type CourseDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDeleteLogic {
	return &CourseDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseDelete 删除课程（透传 RPC）。
func (l *CourseDeleteLogic) CourseDelete(req *types.IdPathReq) (resp *types.NameExistVO, err error) {
	_, gerr := l.svcCtx.CourseRpc.CourseDelete(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "删除课程失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
