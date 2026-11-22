// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package recommend

import (
	"context"

	backendlogic "tjxt/app/search/api/backend/logic"
	searchpb "tjxt/app/search/api/backend/pb"
	"tjxt/app/search/api/internal/svc"
	"tjxt/app/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTopCoursesByCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTopCoursesByCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTopCoursesByCategoryLogic {
	return &GetTopCoursesByCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetTopCoursesByCategory 根据二级分类id查询课程TOP10。
func (l *GetTopCoursesByCategoryLogic) GetTopCoursesByCategory(req *types.IdPathReq) (resp *types.CourseListVO, err error) {
	reply, err := backendlogic.NewGetTopCoursesByCategoryLogic(l.ctx, l.svcCtx.Backend).GetTopCoursesByCategory(&searchpb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	list := make([]types.CourseVO, 0, len(reply.Items))
	for _, item := range reply.Items {
		list = append(list, toCourseVO(item))
	}
	return &types.CourseListVO{List: list}, nil
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
