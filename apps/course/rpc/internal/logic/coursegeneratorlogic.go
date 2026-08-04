package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseGeneratorLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseGeneratorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGeneratorLogic {
	return &CourseGeneratorLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseGeneratorLogic) CourseGenerator(in *pb.Empty) (*pb.IdResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.IdResponse{}, nil
}
