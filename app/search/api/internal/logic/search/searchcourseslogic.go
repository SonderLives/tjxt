// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package search

import (
	"context"

	backendlogic "tjxt/app/search/api/backend/logic"
	searchpb "tjxt/app/search/api/backend/pb"
	"tjxt/app/search/api/internal/svc"
	"tjxt/app/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCoursesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCoursesLogic {
	return &SearchCoursesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SearchCourses 课程全文检索。
func (l *SearchCoursesLogic) SearchCourses(req *types.CourseSearchReq) (resp *types.CourseSearchPageVO, err error) {
	reply, err := backendlogic.NewSearchCoursesLogic(l.ctx, l.svcCtx.Backend).SearchCourses(&searchpb.CourseSearchRequest{
		Keyword:  req.Keyword,
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Sort:     req.Sort,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.CourseVO, 0, len(reply.Items))
	for _, item := range reply.Items {
		list = append(list, toCourseVO(item))
	}
	return &types.CourseSearchPageVO{Total: reply.Total, List: list}, nil
}

func toCourseVO(item *searchpb.CourseVO) types.CourseVO {
	return types.CourseVO{
		Id:            item.Id,
		Name:          item.Name,
		CoverUrl:      item.CoverUrl,
		Price:         item.Price,
		Score:         item.Score,
		Sold:          item.Sold,
		Sections:      item.Sections,
		Free:          item.Free,
		CourseType:    item.CourseType,
		Enable:        item.Enable,
		CategoryIdLv1: item.CategoryIdLv1,
		CategoryIdLv2: item.CategoryIdLv2,
		CategoryIdLv3: item.CategoryIdLv3,
		CreateTime:    item.CreateTime,
		PublishTime:   item.PublishTime,
		Duration:      item.Duration,
	}
}
