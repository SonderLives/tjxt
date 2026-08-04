package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogueIndexListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogueIndexListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogueIndexListLogic {
	return &CourseCatalogueIndexListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCatalogueIndexListLogic) CourseCatalogueIndexList(in *pb.IdRequest) (*pb.CataSimpleList, error) {
	// todo: add your logic here and delete this line

	return &pb.CataSimpleList{}, nil
}
