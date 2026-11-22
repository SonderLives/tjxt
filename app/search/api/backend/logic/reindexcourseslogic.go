package searchlogic

import (
	"context"

	"tjxt/app/search/api/backend/svc"
	"tjxt/app/search/api/backend/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReindexCoursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReindexCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReindexCoursesLogic {
	return &ReindexCoursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReindexCourses 手动触发全量重建 ES 索引（从 course 服务拉取全部已上架课程写入 ES）。
func (l *ReindexCoursesLogic) ReindexCourses(in *pb.Empty) (*pb.ReindexReply, error) {
	indexed, total, err := svc.ReindexAll(l.ctx, l.svcCtx)
	if err != nil {
		return nil, err
	}
	return &pb.ReindexReply{Indexed: indexed, Total: total}, nil
}
