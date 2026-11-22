// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"

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
		return nil, err
	}
	return &types.NameExistVO{Existed: false}, nil
}
