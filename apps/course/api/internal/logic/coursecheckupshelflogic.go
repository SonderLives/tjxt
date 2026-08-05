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

type CourseCheckUpShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCheckUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckUpShelfLogic {
	return &CourseCheckUpShelfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseCheckUpShelf 校验课程是否满足上架条件（透传 RPC）。
func (l *CourseCheckUpShelfLogic) CourseCheckUpShelf(req *types.IdPathReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CourseCheckUpShelf(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "校验课程上架条件失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
