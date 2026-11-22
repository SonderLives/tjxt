// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	backendlogic "tjxt/app/exam/api/backend/logic"
	"tjxt/app/exam/api/backend/pb"
	"tjxt/app/exam/api/internal/svc"
	"tjxt/app/exam/api/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveQuestionBizLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveQuestionBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveQuestionBizLogic {
	return &RemoveQuestionBizLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RemoveQuestionBizLogic) RemoveQuestionBiz(req *types.QuestionBizReq) (resp *types.OkVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	if _, err := backendlogic.NewRemoveQuestionBizLogic(l.ctx, l.svcCtx.Backend).RemoveQuestionBiz(&pb.QuestionBizReq{
		BizId:      req.BizId,
		QuestionId: req.QuestionId,
	}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
