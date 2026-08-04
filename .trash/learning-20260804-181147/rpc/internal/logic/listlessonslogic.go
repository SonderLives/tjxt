package logic

import (
	"context"

	"tjxt/apps/learning/rpc/internal/svc"
	"tjxt/apps/learning/rpc/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLessonsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLessonsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLessonsLogic {
	return &ListLessonsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLessonsLogic) ListLessons(in *pb.ListLessonsRequest) (*pb.ListLessonsReply, error) {
	pageNo, pageSize := in.PageNo, in.PageSize
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	lessons, total, err := l.svcCtx.LessonService.ListLessons(l.ctx, in.UserId, pageNo, pageSize, in.IsAsc)
	if err != nil {
		return nil, err
	}
	reply := &pb.ListLessonsReply{Total: total, Pages: (total + pageSize - 1) / pageSize}
	for index := range lessons {
		reply.List = append(reply.List, lessonReply(lessons[index]))
	}
	return reply, nil
}
