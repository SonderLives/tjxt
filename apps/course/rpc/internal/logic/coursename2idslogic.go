package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseName2IdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseName2IdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseName2IdsLogic {
	return &CourseName2IdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseName2IdsLogic) CourseName2Ids(in *pb.CourseNameRequest) (*pb.CourseIdList, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseIdList{}, nil
}
