package logic

import (
	"context"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddQuestionBizLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddQuestionBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddQuestionBizLogic {
	return &AddQuestionBizLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 题目业务关联
func (l *AddQuestionBizLogic) AddQuestionBiz(in *pb.QuestionBizReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
