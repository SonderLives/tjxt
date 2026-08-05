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

type CourseCataIndexListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCataIndexListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCataIndexListLogic {
	return &CourseCataIndexListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseCataIndexList 查询课程目录的简要序号列表。
func (l *CourseCataIndexListLogic) CourseCataIndexList(req *types.IdPathReq) (resp []types.CataSimpleInfoVO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CourseCatalogueIndexList(l.ctx, &pb.IdRequest{Id: req.Id})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程目录序号列表失败")
	}
	resp = make([]types.CataSimpleInfoVO, 0, len(list.Items))
	for _, c := range list.Items {
		resp = append(resp, cataIndexToSimpleInfoVO(c))
	}
	return resp, nil
}

// cataIndexToSimpleInfoVO pb.CataSimple -> API CataSimpleInfoVO。
func cataIndexToSimpleInfoVO(c *pb.CataSimple) types.CataSimpleInfoVO {
	if c == nil {
		return types.CataSimpleInfoVO{}
	}
	return types.CataSimpleInfoVO{
		Id:           c.Id,
		Name:         c.Name,
		Index:        int64(c.Index),
		ChapterIndex: int64(c.ChapterIndex),
		CIndex:       int64(c.CIndex),
	}
}
