package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSmsPlatformsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSmsPlatformsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSmsPlatformsLogic {
	return &ListSmsPlatformsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短信平台 全量列表，按 priority 升序（数字越小优先级越高）
func (l *ListSmsPlatformsLogic) ListSmsPlatforms(in *pb.Empty) (*pb.SmsPlatformListReply, error) {
	list, err := l.svcCtx.SmsThirdPlatformModel.FindAll(l.ctx)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询短信平台列表失败")
	}
	resp := &pb.SmsPlatformListReply{
		List: make([]*pb.SmsPlatformVO, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, &pb.SmsPlatformVO{
			Id:       item.Id,
			Name:     item.Name,
			Code:     item.Code,
			Priority: int32(item.Priority),
			Status:   int32(item.Status),
		})
	}
	return resp, nil
}
