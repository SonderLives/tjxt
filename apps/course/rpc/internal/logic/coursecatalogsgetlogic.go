package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogsGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCatalogsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogsGetLogic {
	return &CourseCatalogsGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程目录 (catalogue) =====
func (l *CourseCatalogsGetLogic) CourseCatalogsGet(in *pb.IdRequest) (*pb.CourseAndSectionView, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseAndSectionView{}, nil
}
