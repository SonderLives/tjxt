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

type CatalogueBatchQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCatalogueBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CatalogueBatchQueryLogic {
	return &CatalogueBatchQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CatalogueBatchQuery 根据目录 id 批量查询目录简要信息。
func (l *CatalogueBatchQueryLogic) CatalogueBatchQuery(req *types.CourseBatchQueryReq) (resp []types.CataSimpleInfoVO, err error) {
	ids := parseIds(req.Ids)
	if len(ids) == 0 {
		return []types.CataSimpleInfoVO{}, nil
	}
	list, gerr := l.svcCtx.CourseRpc.CourseCatalogueBatchQuery(l.ctx, &pb.IdsRequest{Ids: ids})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "批量查询课程目录失败")
	}
	resp = make([]types.CataSimpleInfoVO, 0, len(list.Items))
	for _, c := range list.Items {
		if c == nil {
			continue
		}
		resp = append(resp, types.CataSimpleInfoVO{
			Id:           c.Id,
			Name:         c.Name,
			Index:        int64(c.Index),
			ChapterIndex: int64(c.ChapterIndex),
			CIndex:       int64(c.CIndex),
		})
	}
	return resp, nil
}
