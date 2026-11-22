package messagelogic

import (
	"context"
	"errors"

	"tjxt/service/message/api/backend/model"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/backend/svc"
	"tjxt/common/auth"
	"tjxt/common/errorx"

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
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	data, err := l.svcCtx.UserInboxModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("站内信不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询站内信失败")
	}
	if data.UserId != userId {
		return nil, errorx.Forbidden("无权操作该站内信")
	}
	if data.IsRead == 1 {
		return &pb.Empty{}, nil
	}
	data.IsRead = 1
	if err := l.svcCtx.UserInboxModel.Update(l.ctx, data); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "标记站内信已读失败")
	}
	return &pb.Empty{}, nil
}
