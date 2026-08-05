package remarklogic

import (
	"context"

	"tjxt/apps/remark/rpc/internal/svc"
	"tjxt/apps/remark/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryLikedBizIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryLikedBizIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryLikedBizIdsLogic {
	return &QueryLikedBizIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *QueryLikedBizIdsLogic) QueryLikedBizIds(in *pb.LikedReq) (*pb.LikedResp, error) {
	ids, err := l.svcCtx.LikeRecordModel.FindLikedBizIds(l.ctx, in.UserId, in.BizType, in.BizIds)
	if err != nil {
		return nil, err
	}
	return &pb.LikedResp{LikedBizIds: ids}, nil
}