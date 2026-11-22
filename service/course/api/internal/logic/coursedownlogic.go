// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/common/errorx"
	"tjxt/service/course/rpc/pb"

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
		return nil, errorx.BadRequestf("课程id不能为空")
	}
	_, err = l.svcCtx.CourseRpc.CourseDown(l.ctx, &pb.IdsRequest{Ids: ids})
	if err != nil {
		return nil, err
	}
	return &types.NameExistVO{Existed: false}, nil
}
