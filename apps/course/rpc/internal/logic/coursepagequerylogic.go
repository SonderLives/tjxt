package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoursePageQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCoursePageQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePageQueryLogic {
	return &CoursePageQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CoursePageQueryLogic) CoursePageQuery(in *pb.CoursePageQueryRequest) (*pb.CoursePageQueryReply, error) {
	// todo: add your logic here and delete this line

	return &pb.CoursePageQueryReply{}, nil
}
