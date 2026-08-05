package logic

import (
	"context"
	"time"

	"tjxt/apps/message/rpc/internal/model"
	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNoticeLogic {
	return &SendNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 站内信 直接发送一条（不查模板）
func (l *SendNoticeLogic) SendNotice(in *pb.SendNoticeReq) (*pb.IdReply, error) {
	if in.UserId <= 0 {
		return nil, xerr.BadRequestf("目标用户不能为空")
	}
	if in.Content == "" {
		return nil, xerr.BadRequestf("通知内容不能为空")
	}

	now := time.Now()
	expireTime := now.AddDate(0, 0, 30)
	if in.ExpireTime != "" {
		t, err := parseDateTime(in.ExpireTime)
		if err != nil {
			return nil, xerr.BadRequestf("expireTime 格式非法")
		}
		expireTime = t.Time
	}

	data := &model.UserInbox{
		Id:         nextID(),
		UserId:     in.UserId,
		Type:       int64(in.Type),
		Title:      in.Title,
		Content:    in.Content,
		IsRead:     0,
		Publisher:  in.Publisher,
		PushTime:   now,
		ExpireTime: expireTime,
	}
	if _, err := l.svcCtx.UserInboxModel.Insert(l.ctx, data); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "发送站内信失败")
	}
	return &pb.IdReply{Id: data.Id}, nil
}
