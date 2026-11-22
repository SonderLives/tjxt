// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSectionGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSectionGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSectionGetLogic {
	return &CourseSectionGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseSectionGet 查询小节信息（透传 RPC）。
func (l *CourseSectionGetLogic) CourseSectionGet(req *types.IdPathReq) (resp *types.CourseSectionInfoVO, err error) {
	info, gerr := l.svcCtx.CourseRpc.CourseSectionGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, errorx.Wrap(gerr, errorx.CodeInternal, "查询小节信息失败")
	}
	return &types.CourseSectionInfoVO{
		Id:       info.Id,
		CourseId: info.CourseId,
		MediaId:  info.MediaId,
	}, nil
}
