package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueSaveLogic {
	return &CourseCatalogueSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCatalogueSaveLogic) CourseCatalogueSave(in *pb.CourseCatalogueSaveRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
