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

// 开通课程（来自 trade 的内部 RPC，user_id 取自请求而非 JWT）
func (l *GrantCoursesLogic) GrantCourses(in *pb.GrantCoursesRequest) (*pb.Empty, error) {
	if err := l.svcCtx.LearningService.GrantCourses(l.ctx, in.UserId, in.CourseIds); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
