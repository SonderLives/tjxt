package remarklogic

import (
	"context"

	"tjxt/apps/remark/rpc/internal/model"
	"tjxt/apps/remark/rpc/internal/svc"
	"tjxt/apps/remark/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type LikeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LikeLogic) Like(in *pb.LikeReq) (*pb.Empty, error) {
	liked := int64(0)
	if in.Liked {
		liked = 1
	}

	m := l.svcCtx.LikeRecordModel
	existing, err := m.FindOneByUserIdBizIdBizType(l.ctx, in.UserId, in.BizId, in.BizType)
	switch {
	case err == sqlc.ErrNotFound:
		_, err = m.Insert(l.ctx, &model.LikeRecord{
			UserId:  in.UserId,
			BizId:   in.BizId,
			BizType: in.BizType,
			Liked:   liked,
		})
		return &pb.Empty{}, err
	case err != nil:
		return nil, err
	}
	if existing.Liked == liked {
		return &pb.Empty{}, nil
	}
	existing.Liked = liked
	return &pb.Empty{}, m.Update(l.ctx, existing)
}