// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"context"

	backendlogic "tjxt/service/exam/api/backend/logic"
	"tjxt/service/exam/api/backend/pb"
	"tjxt/service/exam/api/internal/logic"
	"tjxt/service/exam/api/internal/svc"
	"tjxt/service/exam/api/internal/types"
	"tjxt/common/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionLogic {
	return &GetQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionLogic) GetQuestion(req *types.IdPathReq) (resp *types.QuestionVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := backendlogic.NewGetQuestionLogic(l.ctx, l.svcCtx.Backend).GetQuestion(&pb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return logic.ToQuestionVO(rpcResp), nil
}
