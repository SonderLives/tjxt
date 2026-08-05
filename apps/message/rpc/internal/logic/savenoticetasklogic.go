package logic

import (
	"context"
	"database/sql"

	"tjxt/apps/message/rpc/internal/model"
	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoticeTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoticeTaskLogic {
	return &SaveNoticeTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知任务 新增/更新
func (l *SaveNoticeTaskLogic) SaveNoticeTask(in *pb.NoticeTaskSaveReq) (*pb.IdReply, error) {
	if in.Name == "" {
		return nil, xerr.BadRequestf("任务名称不能为空")
	}

	// 校验通知模板存在
	if _, err := l.svcCtx.NoticeTemplateModel.FindOne(l.ctx, in.TemplateId); err != nil {
		if isNotFound(err) {
			return nil, xerr.BadRequestf("通知模板不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询通知模板失败")
	}

	pushTime, err := parseDateTime(in.PushTime)
	if err != nil {
		return nil, xerr.BadRequestf("pushTime 格式非法")
	}
	expireTime, err := parseDateTime(in.ExpireTime)
	if err != nil {
		return nil, xerr.BadRequestf("expireTime 格式非法")
	}
	interval := sql.NullInt64{}
	if in.Interval > 0 {
		interval = sql.NullInt64{Int64: int64(in.Interval), Valid: true}
	}
	maxTimes := in.MaxTimes
	if maxTimes <= 0 {
		maxTimes = 1
	}
	operator := userIdFromCtx(l.ctx)
	partial := byte(0)
	if in.Partial {
		partial = 1
	}

	if in.Id == 0 {
		data := &model.NoticeTask{
			Id:         nextID(),
			TemplateId: in.TemplateId,
			Name:       in.Name,
			Partial:    partial,
			PushTime:   pushTime,
			Interval:   interval,
			ExpireTime: expireTime,
			MaxTimes:   int64(maxTimes),
			Finished:   0,
			Creater:    operator,
			Updater:    operator,
		}
		if _, err := l.svcCtx.NoticeTaskModel.Insert(l.ctx, data); err != nil {
			return nil, xerr.Wrapf(err, xerr.CodeInternal, "新增通知任务失败")
		}
		return &pb.IdReply{Id: data.Id}, nil
	}

	old, err := l.svcCtx.NoticeTaskModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("通知任务不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询通知任务失败")
	}
	old.TemplateId = in.TemplateId
	old.Name = in.Name
	old.Partial = partial
	old.PushTime = pushTime
	old.Interval = interval
	old.ExpireTime = expireTime
	old.MaxTimes = int64(maxTimes)
	old.Updater = operator
	if err := l.svcCtx.NoticeTaskModel.Update(l.ctx, old); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新通知任务失败")
	}
	return &pb.IdReply{Id: old.Id}, nil
}
