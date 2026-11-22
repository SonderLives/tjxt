// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"context"

	backendlogic "tjxt/api/exam/backend/logic"
	"tjxt/api/exam/backend/pb"
	"tjxt/api/exam/internal/logic"
	"tjxt/api/exam/internal/svc"
	"tjxt/api/exam/internal/types"
	"tjxt/pkg/auth"

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
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := backendlogic.NewGetQuestionsByBizLogic(l.ctx, l.svcCtx.Backend).GetQuestionsByBiz(&pb.QuestionBizListReq{
		BizId:    req.BizId,
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
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
