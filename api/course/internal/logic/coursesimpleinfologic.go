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

type CourseSimpleInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSimpleInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSimpleInfoLogic {
	return &CourseSimpleInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseSimpleInfo 批量查询课程简单信息（透传 RPC）。
func (l *CourseSimpleInfoLogic) CourseSimpleInfo(req *types.CourseSimpleInfoQueryReq) (resp []types.CourseSimpleInfoDTO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CourseSimpleInfoList(l.ctx, &pb.CourseSimpleInfoQueryRequest{
		Ids:          parseIds(req.Ids),
		ThirdCateIds: parseIds(req.ThirdCataIds),
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程简单信息失败")
	}
	resp = make([]types.CourseSimpleInfoDTO, 0, len(list.Items))
	for _, item := range list.Items {
		resp = append(resp, types.CourseSimpleInfoDTO{
			Id:              item.Id,
			Name:            item.Name,
			CoverUrl:        item.CoverUrl,
			Price:           item.Price,
			Free:            int64(item.Free),
			SectionNum:      item.SectionNum,
			Status:          int64(item.Status),
			ValidDuration:   item.ValidDuration,
			PurchaseEndTime: item.PurchaseEndTime,
			FirstCateId:     item.FirstCateId,
			SecondCateId:    item.SecondCateId,
			ThirdCateId:     item.ThirdCateId,
		})
	}
	return resp, nil
}
