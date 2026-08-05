// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"context"

	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"
	examclient "tjxt/apps/exam/rpc/exam"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveQuestionLogic {
	return &SaveQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveQuestionLogic) SaveQuestion(req *types.QuestionSaveReq) (resp *types.IdVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ExamRpc.SaveQuestion(l.ctx, &examclient.QuestionSaveReq{
		Id:         req.Id,
		Name:       req.Name,
		Type:       req.Type,
		CateId1:    req.CateId1,
		CateId2:    req.CateId2,
		CateId3:    req.CateId3,
		Difficulty: req.Difficulty,
		Score:      req.Score,
		Options:    req.Options,
		Answer:     req.Answer,
		Analysis:   req.Analysis,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: rpcResp.Id}, nil
}
