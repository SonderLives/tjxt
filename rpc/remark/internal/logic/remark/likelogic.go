package remarklogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/rpc/remark/internal/model"
	"tjxt/rpc/remark/internal/svc"
	"tjxt/rpc/remark/pb"

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
	case errors.Is(err, sqlc.ErrNotFound):
		if _, err = m.Insert(l.ctx, &model.LikeRecord{
			UserId:  in.UserId,
			BizId:   in.BizId,
			BizType: in.BizType,
			Liked:   liked,
		}); err != nil {
			return nil, errorx.Wrap(err, errorx.CodeInternal, "点赞失败")
		}
		return &pb.Empty{}, nil
	case err != nil:
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询点赞状态失败")
	}
	if existing.Liked == liked {
		return &pb.Empty{}, nil
	}
	existing.Liked = liked
	if err := m.Update(l.ctx, existing); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "更新点赞状态失败")
	}
	return &pb.Empty{}, nil
}
