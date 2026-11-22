package like

import (
	"context"

	backendlogic "tjxt/api/remark/backend/logic"
	"tjxt/api/remark/backend/pb"
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
	_, err = backendlogic.NewLikeLogic(l.ctx, l.svcCtx.Backend).Like(&pb.LikeReq{
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
