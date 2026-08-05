package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/model"
	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePublicNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSavePublicNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePublicNoticeLogic {
	return &SavePublicNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 公告 新增/更新
func (l *SavePublicNoticeLogic) SavePublicNotice(in *pb.PublicNoticeSaveReq) (*pb.IdReply, error) {
	if in.Title == "" || in.Content == "" {
		return nil, xerr.BadRequestf("公告标题/内容不能为空")
	}
	pushTime, err := parseDateTime(in.PushTime)
	if err != nil {
		return nil, xerr.BadRequestf("pushTime 必填且格式非法")
	}
	expireTime, err := parseDateTime(in.ExpireTime)
	if err != nil {
		return nil, xerr.BadRequestf("expireTime 必填且格式非法")
	}
	if !pushTime.Valid || !expireTime.Valid {
		return nil, xerr.BadRequestf("pushTime/expireTime 必填")
	}

	if in.Id == 0 {
		data := &model.PublicNotice{
			Id:         nextID(),
			Title:      in.Title,
			Content:    in.Content,
			Type:       int64(in.Type),
			PushTime:   pushTime.Time,
			ExpireTime: expireTime.Time,
		}
		if _, err := l.svcCtx.PublicNoticeModel.Insert(l.ctx, data); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "新增公告失败")
		}
		return &pb.IdReply{Id: data.Id}, nil
	}

	old, err := l.svcCtx.PublicNoticeModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("公告不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询公告失败")
	}
	old.Title = in.Title
	old.Content = in.Content
	old.Type = int64(in.Type)
	old.PushTime = pushTime.Time
	old.ExpireTime = expireTime.Time
	if err := l.svcCtx.PublicNoticeModel.Update(l.ctx, old); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新公告失败")
	}
	return &pb.IdReply{Id: old.Id}, nil
}
