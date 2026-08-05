// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
