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

type CourseGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGetLogic {
	return &CourseGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CourseGet 查询课程完整信息（透传 RPC）。
func (l *CourseGetLogic) CourseGet(req *types.CourseFullInfoQueryReq) (resp *types.CourseFullInfoDTO, err error) {
	info, gerr := l.svcCtx.CourseRpc.CourseFullInfoGet(l.ctx, &pb.CourseFullInfoGetRequest{
		Id:            req.Id,
		WithCatalogue: req.WithCatalogue,
		WithTeachers:  req.WithTeachers,
	})
	if gerr != nil {
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询课程详情失败")
	}
	return &types.CourseFullInfoDTO{
		Id:              info.Id,
		Name:            info.Name,
		CoverUrl:        info.CoverUrl,
		Price:           info.Price,
		FirstCateId:     info.FirstCateId,
		SecondCateId:    info.SecondCateId,
		ThirdCateId:     info.ThirdCateId,
		ValidDuration:   info.ValidDuration,
		PurchaseEndTime: info.PurchaseEndTime,
		SectionNum:      info.SectionNum,
		TeacherIds:      info.TeacherIds,
		Chapters:        toFullInfoCatalogueVOs(info.Chapters),
	}, nil
}

// toFullInfoCatalogueVOs pb.CourseChapterInfo 列表 -> API CatalogueVO 列表（递归）。
func toFullInfoCatalogueVOs(list []*pb.CourseChapterInfo) []types.CatalogueVO {
	if len(list) == 0 {
		return nil
	}
	vos := make([]types.CatalogueVO, 0, len(list))
	for _, c := range list {
		vos = append(vos, types.CatalogueVO{
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
			Sections:      toFullInfoCatalogueVOs(c.Sections),
			CanUpdate:     c.CanUpdate,
		})
	}
	return vos
}
