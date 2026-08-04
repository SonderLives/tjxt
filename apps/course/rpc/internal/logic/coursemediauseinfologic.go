package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseMediaUseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCourseMediaUseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaUseInfoLogic {
	return &CourseMediaUseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CourseMediaUseInfoLogic) CourseMediaUseInfo(in *pb.MediaIdsRequest) (*pb.MediaQuoteList, error) {
	// todo: add your logic here and delete this line

	return &pb.MediaQuoteList{}, nil
}
