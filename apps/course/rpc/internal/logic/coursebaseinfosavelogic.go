package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseBaseInfoSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseBaseInfoSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoSaveLogic {
	return &CourseBaseInfoSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseBaseInfoSaveLogic) CourseBaseInfoSave(in *pb.CourseBaseInfoSaveRequest) (*pb.IdResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.IdResponse{}, nil
}
