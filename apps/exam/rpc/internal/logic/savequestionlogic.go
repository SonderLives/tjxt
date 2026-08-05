package logic

import (
	"context"

	"tjxt/apps/exam/rpc/internal/svc"
	"tjxt/apps/exam/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveQuestionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveQuestionLogic {
	return &SaveQuestionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 题目管理
func (l *SaveQuestionLogic) SaveQuestion(in *pb.QuestionSaveReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
