package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSectionGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSectionGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSectionGetLogic {
	return &CourseSectionGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseSectionGet 按小节 id 查询小节所属课程与媒资信息（正式目录表）。
func (l *CourseSectionGetLogic) CourseSectionGet(in *pb.IdRequest) (*pb.CourseSectionInfo, error) {
	if in.Id == 0 {
		return nil, xerr.BadRequestf("小节id不能为空")
	}
	c, err := l.svcCtx.CourseCatalogueModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("小节不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询小节失败")
	}
	return &pb.CourseSectionInfo{
		Id:       c.Id,
		CourseId: c.CourseId,
		MediaId:  c.MediaId,
	}, nil
}
