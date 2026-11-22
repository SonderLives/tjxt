// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	backendlogic "tjxt/api/exam/backend/logic"
	"tjxt/api/exam/backend/pb"
	"tjxt/api/exam/internal/svc"
	"tjxt/api/exam/internal/types"
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
	rpcResp, err := backendlogic.NewAddQuestionBizLogic(l.ctx, l.svcCtx.Backend).AddQuestionBiz(&pb.QuestionBizReq{
		BizId:      req.BizId,
		QuestionId: req.QuestionId,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: rpcResp.Id}, nil
}
