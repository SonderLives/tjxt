// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
