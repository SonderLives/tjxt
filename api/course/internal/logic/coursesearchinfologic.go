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

type CourseSearchInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSearchInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSearchInfoLogic {
	return &CourseSearchInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseSearchInfo 查询课程搜索索引信息（透传 RPC）。
func (l *CourseSearchInfoLogic) CourseSearchInfo(req *types.IdPathReq) (resp *types.CourseSearchIndexVO, err error) {
	info, gerr := l.svcCtx.CourseRpc.CourseSearchInfoForIndex(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程搜索信息失败")
	}
	return &types.CourseSearchIndexVO{
		Id:            info.Id,
		Name:          info.Name,
		CoverUrl:      info.CoverUrl,
		Price:         info.Price,
		Score:         info.Score,
		Sold:          info.Sold,
		Sections:      info.Sections,
		Free:          int64(info.Free),
		CourseType:    int64(info.CourseType),
		Enable:        int64(info.Enable),
		CategoryIdLv1: info.CategoryIdLv1,
		CategoryIdLv2: info.CategoryIdLv2,
		CategoryIdLv3: info.CategoryIdLv3,
		CreateTime:    info.CreateTime,
		PublishTime:   info.PublishTime,
		Duration:      info.Duration,
	}, nil
}
