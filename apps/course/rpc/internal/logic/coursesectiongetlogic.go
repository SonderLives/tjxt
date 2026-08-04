package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSectionGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseSectionGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSectionGetLogic {
	return &CourseSectionGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseSectionGetLogic) CourseSectionGet(in *pb.IdRequest) (*pb.CourseSectionInfo, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseSectionInfo{}, nil
}
