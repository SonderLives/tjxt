package like

import (
	"context"
	"strconv"

	backendlogic "tjxt/api/remark/backend/logic"
	"tjxt/api/remark/backend/pb"
	"tjxt/pkg/auth"

	"tjxt/api/remark/internal/svc"
	"tjxt/api/remark/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLikedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLikedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLikedLogic {
	return &ListLikedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLikedLogic) ListLiked(req *types.LikeListReq) (resp *types.LikeListResp, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	bizIds := make([]int64, 0, len(req.BizIds))
	for _, s := range req.BizIds {
		id, perr := strconv.ParseInt(s, 10, 64)
		if perr == nil && id > 0 {
			bizIds = append(bizIds, id)
		}
	}
	r, err := backendlogic.NewQueryLikedBizIdsLogic(l.ctx, l.svcCtx.Backend).QueryLikedBizIds(&pb.LikedReq{
		UserId:  userId,
		BizType: req.BizType,
		BizIds:  bizIds,
	})
	if err != nil {
		return nil, err
	}
	return &types.LikeListResp{LikedBizIds: r.LikedBizIds}, nil
}
