package remarklogic

import (
	"context"

	"tjxt/pkg/errorx"
	"tjxt/rpc/remark/internal/svc"
	"tjxt/rpc/remark/pb"

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
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询点赞状态失败")
	}
	return &pb.LikedResp{LikedBizIds: ids}, nil
}
