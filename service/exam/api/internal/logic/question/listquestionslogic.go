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
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := backendlogic.NewListQuestionsLogic(l.ctx, l.svcCtx.Backend).ListQuestions(&pb.QuestionListReq{
		PageNo:     int32(req.PageNo),
		PageSize:   int32(req.PageSize),
		Name:       req.Name,
		Type:       req.Type,
		CateId1:    req.CateId1,
		CateId2:    req.CateId2,
		Difficulty: req.Difficulty,
	})
	if err != nil {
		return nil, err
	}
	vo := &types.QuestionListVO{
		Total: rpcResp.Total,
		List:  make([]types.QuestionVO, 0, len(rpcResp.List)),
	}
	for _, item := range rpcResp.List {
		vo.List = append(vo.List, *logic.ToQuestionVO(item))
	}
	return vo, nil
}
