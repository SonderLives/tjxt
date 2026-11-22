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

type CourseCheckNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCheckNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckNameLogic {
	return &CourseCheckNameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseCheckName 校验课程名称是否重复（透传 RPC）。
func (l *CourseCheckNameLogic) CourseCheckName(req *types.CourseCheckNameReq) (resp *types.NameExistVO, err error) {
	res, gerr := l.svcCtx.CourseRpc.CourseCheckName(l.ctx, &pb.CourseCheckNameRequest{
		Name: req.Name,
		Id:   req.Id,
	})
	if gerr != nil {
		return nil, gerr
	}
	return &types.NameExistVO{Existed: res.Existed}, nil
}
