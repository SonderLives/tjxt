package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueTreeGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueTreeGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueTreeGetLogic {
	return &CourseCatalogueTreeGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCatalogueTreeGetLogic) CourseCatalogueTreeGet(in *pb.CourseCatalogueQueryRequest) (*pb.CatalogueTreeList, error) {
	// todo: add your logic here and delete this line

	return &pb.CatalogueTreeList{}, nil
}
