package courselogic

import (
	"context"

	"tjxt/rpc/course/internal/svc"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.BadRequestf("小节id不能为空")
	}
	c, err := l.svcCtx.CourseCatalogueModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("小节不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询小节失败")
	}
	return &pb.CourseSectionInfo{
		Id:       c.Id,
		CourseId: c.CourseId,
		MediaId:  c.MediaId,
	}, nil
}
