package messagelogic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"
	"tjxt/rpc/message/internal/svc"
	"tjxt/rpc/message/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteInboxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteInboxLogic {
	return &DeleteInboxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 站内信 物理删除（校验归属当前用户）
func (l *DeleteInboxLogic) DeleteInbox(in *pb.IdReq) (*pb.Empty, error) {
	userId, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	data, err := l.svcCtx.UserInboxModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("站内信不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询站内信失败")
	}
	if data.UserId != userId {
		return nil, errorx.Forbidden("无权操作该站内信")
	}
	if err := l.svcCtx.UserInboxModel.Delete(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "删除站内信失败")
	}
	return &pb.Empty{}, nil
}
