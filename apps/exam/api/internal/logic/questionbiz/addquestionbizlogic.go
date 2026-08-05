// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"
	examclient "tjxt/apps/exam/rpc/exam"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddQuestionBizLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddQuestionBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddQuestionBizLogic {
	return &AddQuestionBizLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddQuestionBizLogic) AddQuestionBiz(req *types.QuestionBizReq) (resp *types.IdVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.ExamRpc.AddQuestionBiz(l.ctx, &examclient.QuestionBizReq{
		BizId:      req.BizId,
		QuestionId: req.QuestionId,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: rpcResp.Id}, nil
}
