package courselogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/rpc/course/internal/model"
	"tjxt/rpc/course/internal/svc"
	"tjxt/rpc/course/pb"

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
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("小节不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询小节失败")
	}
	return &pb.CourseSectionInfo{
		Id:       cata.Id,
		CourseId: cata.CourseId,
		MediaId:  cata.MediaId,
	}, nil
}
