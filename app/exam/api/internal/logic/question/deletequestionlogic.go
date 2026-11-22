// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"context"

	backendlogic "tjxt/app/exam/api/backend/logic"
	"tjxt/app/exam/api/backend/pb"
	"tjxt/app/exam/api/internal/svc"
	"tjxt/app/exam/api/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteQuestionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteQuestionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteQuestionLogic {
	return &DeleteQuestionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteQuestionLogic) DeleteQuestion(req *types.IdPathReq) (resp *types.OkVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	if _, err := backendlogic.NewDeleteQuestionLogic(l.ctx, l.svcCtx.Backend).DeleteQuestion(&pb.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
