package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseDownLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseDownLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDownLogic {
	return &CourseDownLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseDown 批量下架：逐个把正式课程状态置为已下架。
func (l *CourseDownLogic) CourseDown(in *pb.IdsRequest) (*pb.Empty, error) {
	for _, id := range in.Ids {
		if err := downCourse(l.ctx, l.svcCtx, id); err != nil {
			return nil, err
		}
	}
	return &pb.Empty{}, nil
}
