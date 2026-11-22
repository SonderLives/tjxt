// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package recommend

import (
	"context"

	"tjxt/apps/search/api/internal/svc"
	"tjxt/apps/search/api/internal/types"
	searchclient "tjxt/apps/search/rpc/client/search"
	searchpb "tjxt/apps/search/rpc/pb"

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
	rpcResp, err := l.svcCtx.SearchRpc.GetTopCoursesByCategory(l.ctx, &searchclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}

	list := make([]types.CourseVO, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
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
