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

type CourseDownShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDownShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDownShelfLogic {
	return &CourseDownShelfLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseDownShelf 课程下架（透传 RPC）。
func (l *CourseDownShelfLogic) CourseDownShelf(req *types.IdPathReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CourseDownShelf(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.NameExistVO{Existed: false}, nil
}
