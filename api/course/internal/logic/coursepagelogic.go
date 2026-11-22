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

type CoursePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoursePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePageLogic {
	return &CoursePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CoursePage 管理端课程分页查询（透传 RPC）。
func (l *CoursePageLogic) CoursePage(req *types.CoursePageReq) (resp *types.CoursePageReply, err error) {
	page, gerr := l.svcCtx.CourseRpc.CoursePageQuery(l.ctx, &pb.CoursePageQueryRequest{
		PageNo:       req.PageNo,
		PageSize:     req.PageSize,
		IsAsc:        req.IsAsc,
		SortBy:       req.SortBy,
		Keyword:      req.Keyword,
		Status:       int32(req.Status),
		Free:         int32(req.Free),
		CourseType:   int32(req.CourseType),
		FirstCateId:  req.FirstCateId,
		SecondCateId: req.SecondCateId,
		ThirdCateId:  req.ThirdCateId,
		BeginTime:    req.BeginTime,
		EndTime:      req.EndTime,
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "分页查询课程失败")
	}
	return toCoursePageReply(page), nil
}

// toCoursePageReply pb.CoursePageQueryReply -> API CoursePageReply
func toCoursePageReply(page *pb.CoursePageQueryReply) *types.CoursePageReply {
	resp := &types.CoursePageReply{
		Total: page.Total,
		Pages: page.Pages,
		List:  make([]types.CoursePageVO, 0, len(page.List)),
	}
	for _, item := range page.List {
		resp.List = append(resp.List, types.CoursePageVO{
			Id:             item.Id,
			Name:           item.Name,
			CoverUrl:       item.CoverUrl,
			Price:          item.Price,
			Free:           int64(item.Free),
			Status:         int64(item.Status),
			StatusDesc:     item.StatusDesc,
			FirstCateName:  item.FirstCateName,
			SecondCateName: item.SecondCateName,
			ThirdCateName:  item.ThirdCateName,
			Score:          item.Score,
			Sold:           item.Sold,
			SectionNum:     item.SectionNum,
			MediaDuration:  item.MediaDuration,
			PublishTime:    item.PublishTime,
			CreateTime:     item.CreateTime,
			UpdateTime:     item.UpdateTime,
		})
	}
	return resp
}
