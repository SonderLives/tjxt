package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSearchInfoForIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSearchInfoForIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSearchInfoForIndexLogic {
	return &CourseSearchInfoForIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseSearchInfoForIndexLogic) CourseSearchInfoForIndex(in *pb.IdRequest) (*pb.CourseSearchIndexInfo, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseSearchIndexInfo{}, nil
}
