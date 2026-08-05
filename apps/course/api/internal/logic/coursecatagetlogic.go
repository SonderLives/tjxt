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

type CourseCataGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCataGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCataGetLogic {
	return &CourseCataGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseCataGet 查询课程目录树（章/节/测试）。
func (l *CourseCataGetLogic) CourseCataGet(req *types.CourseCataQueryReq) (resp []types.CatalogueVO, err error) {
	tree, gerr := l.svcCtx.CourseRpc.CourseCatalogueTreeGet(l.ctx, &pb.CourseCatalogueQueryRequest{
		Id:           req.Id,
		See:          req.See,
		WithPractice: req.WithPractice,
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程目录树失败")
	}
	resp = make([]types.CatalogueVO, 0, len(tree.Items))
	for _, c := range tree.Items {
		resp = append(resp, cataGetToCatalogueVO(c))
	}
	return resp, nil
}

// cataGetToCatalogueVO pb.CourseChapterInfo -> API CatalogueVO（递归映射子节）。
func cataGetToCatalogueVO(c *pb.CourseChapterInfo) types.CatalogueVO {
	if c == nil {
		return types.CatalogueVO{}
	}
	vo := types.CatalogueVO{
		Id:            c.Id,
		Name:          c.Name,
		Index:         int64(c.Index),
		Type:          int64(c.Type),
		MediaId:       c.MediaId,
		MediaName:     c.MediaName,
		MediaDuration: c.MediaDuration,
		Trailer:       c.Trailer,
		SubjectNum:    int64(c.SubjectNum),
		TotalScore:    c.TotalScore,
		CanUpdate:     c.CanUpdate,
	}
	for _, s := range c.Sections {
		vo.Sections = append(vo.Sections, cataGetToCatalogueVO(s))
	}
	return vo
}
