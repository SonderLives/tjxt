package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoursePortalQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCoursePortalQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePortalQueryLogic {
	return &CoursePortalQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CoursePortalQueryLogic) CoursePortalQuery(in *pb.CoursePortalQueryRequest) (*pb.CoursePageQueryReply, error) {
	// todo: add your logic here and delete this line

	return &pb.CoursePageQueryReply{}, nil
}
