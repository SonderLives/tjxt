package like

import (
	"context"

	backendlogic "tjxt/app/remark/api/backend/logic"
	"tjxt/app/remark/api/backend/pb"
	"tjxt/pkg/auth"

	"tjxt/app/remark/api/internal/svc"
	"tjxt/app/remark/api/internal/types"

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
