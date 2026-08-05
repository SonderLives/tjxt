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

type CoursePortalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoursePortalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePortalLogic {
	return &CoursePortalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CoursePortal 门户端课程分页查询（透传 RPC）。
func (l *CoursePortalLogic) CoursePortal(req *types.CoursePortalReq) (resp *types.CoursePageReply, err error) {
	page, gerr := l.svcCtx.CourseRpc.CoursePortalQuery(l.ctx, &pb.CoursePortalQueryRequest{
		PageNo:        req.PageNo,
		PageSize:      req.PageSize,
		IsAsc:         req.IsAsc,
		SortBy:        req.SortBy,
		CategoryIdLv1: req.CategoryIdLv1,
		CategoryIdLv2: req.CategoryIdLv2,
		CategoryIdLv3: req.CategoryIdLv3,
		Free:          int32(req.Free),
		Type:          int32(req.Type),
		Keyword:       req.Keyword,
		BeginTime:     req.BeginTime,
		EndTime:       req.EndTime,
		Status:        int32(req.Status),
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询门户课程列表失败")
	}
	return toCoursePageReply(page), nil
}
