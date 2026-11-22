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

type CourseDownLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDownLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDownLogic {
	return &CourseDownLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseDown 批量下架课程（透传 RPC）。
func (l *CourseDownLogic) CourseDown(req *types.CourseIdListReq) (resp *types.NameExistVO, err error) {
	ids := parseIds(req.CourseIds)
	if len(ids) == 0 {
		return nil, xerr.BadRequestf("课程id不能为空")
	}
	_, err = l.svcCtx.CourseRpc.CourseDown(l.ctx, &pb.IdsRequest{Ids: ids})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "批量下架课程失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
