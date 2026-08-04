package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueBatchQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueBatchQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueBatchQueryLogic {
	return &CourseCatalogueBatchQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCatalogueBatchQueryLogic) CourseCatalogueBatchQuery(in *pb.IdsRequest) (*pb.CataSimpleList, error) {
	// todo: add your logic here and delete this line

	return &pb.CataSimpleList{}, nil
}
