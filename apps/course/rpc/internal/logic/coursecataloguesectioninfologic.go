package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueSectionInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueSectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueSectionInfoLogic {
	return &CourseCatalogueSectionInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseCatalogueSectionInfo 查询小节的媒资信息（正式表 course_catalogue）。
func (l *CourseCatalogueSectionInfoLogic) CourseCatalogueSectionInfo(in *pb.IdRequest) (*pb.CourseSectionInfo, error) {
	cata, err := l.svcCtx.CourseCatalogueModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("小节不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询小节失败")
	}
	return &pb.CourseSectionInfo{
		Id:       cata.Id,
		CourseId: cata.CourseId,
		MediaId:  cata.MediaId,
	}, nil
}
