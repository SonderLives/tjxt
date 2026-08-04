package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueSectionInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueSectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueSectionInfoLogic {
	return &CourseCatalogueSectionInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCatalogueSectionInfoLogic) CourseCatalogueSectionInfo(in *pb.IdRequest) (*pb.CourseSectionInfo, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseSectionInfo{}, nil
}
