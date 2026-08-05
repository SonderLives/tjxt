package logic

import (
	"context"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveQuestionBizLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveQuestionBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveQuestionBizLogic {
	return &RemoveQuestionBizLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveQuestionBizLogic) RemoveQuestionBiz(in *pb.QuestionBizReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
