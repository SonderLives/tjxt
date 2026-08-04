package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseBaseInfoGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseBaseInfoGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoGetLogic {
	return &CourseBaseInfoGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 课程基础信息 =====
func (l *CourseBaseInfoGetLogic) CourseBaseInfoGet(in *pb.IdRequest) (*pb.CourseBaseInfoView, error) {
	// todo: add your logic here and delete this line

	return &pb.CourseBaseInfoView{}, nil
}
