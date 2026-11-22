package like

import (
	"context"

	remarkclient "tjxt/rpc/remark/client/remark"
	"tjxt/pkg/auth"

	"tjxt/api/remark/internal/svc"
	"tjxt/api/remark/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeLogic) Like(req *types.LikeRecordFormReq) (resp *types.OkVO, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.RemarkRpc.Like(l.ctx, &remarkclient.LikeReq{
		UserId:  userId,
		BizId:   req.BizId,
		BizType: req.BizType,
		Liked:   req.Liked,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
