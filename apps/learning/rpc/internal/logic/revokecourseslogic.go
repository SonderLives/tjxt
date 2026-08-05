package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeCoursesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeCoursesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeCoursesLogic {
	return &RevokeCoursesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 撤销课程（来自 trade 的内部 RPC，user_id 取自请求而非 JWT）
func (l *RevokeCoursesLogic) RevokeCourses(in *pb.GrantCoursesRequest) (*pb.Empty, error) {
	if err := l.svcCtx.LearningService.RevokeCourses(l.ctx, in.UserId, in.CourseIds); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
