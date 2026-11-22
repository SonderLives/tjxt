// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/app/course/api/internal/svc"
	"tjxt/app/course/api/internal/types"
	"tjxt/app/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseUpShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseUpShelfLogic {
	return &CourseUpShelfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseUpShelf 课程上架/发布（透传 RPC）。
func (l *CourseUpShelfLogic) CourseUpShelf(req *types.IdPathReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CourseUpShelf(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.NameExistVO{Existed: false}, nil
}
