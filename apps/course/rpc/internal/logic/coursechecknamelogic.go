package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCheckNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseCheckNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckNameLogic {
	return &CourseCheckNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseCheckNameLogic) CourseCheckName(in *pb.CourseCheckNameRequest) (*pb.NameExistReply, error) {
	// todo: add your logic here and delete this line

	return &pb.NameExistReply{}, nil
}
