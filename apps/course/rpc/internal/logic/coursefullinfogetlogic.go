package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseFullInfoGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseFullInfoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseFullInfoGetLogic {
	return &CourseFullInfoGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseFullInfoGetLogic) CourseFullInfoGet(in *pb.CourseFullInfoGetRequest) (*pb.CourseFullInfo, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseFullInfo{}, nil
}
