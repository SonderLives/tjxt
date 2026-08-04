package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSimpleInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSimpleInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSimpleInfoListLogic {
	return &CourseSimpleInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseSimpleInfoListLogic) CourseSimpleInfoList(in *pb.CourseSimpleInfoQueryRequest) (*pb.CourseSimpleInfoListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseSimpleInfoListReply{}, nil
}
