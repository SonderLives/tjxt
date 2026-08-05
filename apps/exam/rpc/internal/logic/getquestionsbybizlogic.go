package logic

import (
	"context"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionsByBizLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionsByBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionsByBizLogic {
	return &GetQuestionsByBizLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionsByBizLogic) GetQuestionsByBiz(in *pb.QuestionBizListReq) (*pb.QuestionListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.QuestionListReply{}, nil
}
