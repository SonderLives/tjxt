// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionsByBizLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQuestionsByBizLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionsByBizLogic {
	return &GetQuestionsByBizLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQuestionsByBizLogic) GetQuestionsByBiz(req *types.QuestionBizListReq) (resp *types.QuestionListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
