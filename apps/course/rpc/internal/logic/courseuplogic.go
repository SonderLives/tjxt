package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseUpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseUpLogic {
	return &CourseUpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CourseUp 批量上架：对每个课程 id 执行与 CourseUpShelf 相同的草稿复制发布流程。
func (l *CourseUpLogic) CourseUp(in *pb.IdsRequest) (*pb.Empty, error) {
	for _, id := range in.Ids {
		if err := publishCourse(l.ctx, l.svcCtx, id); err != nil {
			return nil, err
		}
	}
	return &pb.Empty{}, nil
}
