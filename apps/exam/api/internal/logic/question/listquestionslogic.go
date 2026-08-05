// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"context"

	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListQuestionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListQuestionsLogic {
	return &ListQuestionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListQuestionsLogic) ListQuestions(req *types.QuestionListReq) (resp *types.QuestionListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
