package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GrantCoursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGrantCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantCoursesLogic {
	return &GrantCoursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ---- 内部：来自 trade 的 mq 事件 / 内部 RPC 调用 ----
func (l *GrantCoursesLogic) GrantCourses(in *pb.GrantCoursesRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
