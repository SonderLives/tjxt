package messagelogic

import (
	"context"
	"errors"

	"tjxt/api/message/backend/model"
	"tjxt/api/message/backend/pb"
	"tjxt/api/message/backend/svc"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeTaskLogic {
	return &GetNoticeTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知任务 详情
func (l *GetNoticeTaskLogic) GetNoticeTask(in *pb.IdReq) (*pb.NoticeTaskVO, error) {
	data, err := l.svcCtx.NoticeTaskModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("通知任务不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询通知任务失败")
	}
	return toNoticeTaskVO(data), nil
}
