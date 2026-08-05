package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkInboxReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkInboxReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkInboxReadLogic {
	return &MarkInboxReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 站内信 标记已读（校验归属当前用户）
func (l *MarkInboxReadLogic) MarkInboxRead(in *pb.IdReq) (*pb.Empty, error) {
	userId := userIdFromCtx(l.ctx)
	data, err := l.svcCtx.UserInboxModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("站内信不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询站内信失败")
	}
	if data.UserId != userId {
		return nil, xerr.Forbidden("无权操作该站内信")
	}
	if data.IsRead == 1 {
		return &pb.Empty{}, nil
	}
	data.IsRead = 1
	if err := l.svcCtx.UserInboxModel.Update(l.ctx, data); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "标记站内信已读失败")
	}
	return &pb.Empty{}, nil
}
